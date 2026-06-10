//go:build linux

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/nooga/let-go/pkg/vm"
)

const (
	seccompModeFilter            = 2
	seccompRetKillProcess        = 0x80000000
	seccompRetErrno              = 0x00050000
	seccompRetAllow              = 0x7fff0000
	seccompDataOffsetNR          = 0
	seccompDataOffsetArch        = 4
	auditArchX86_64       uint32 = 0xc000003e
	auditArchAArch64      uint32 = 0xc00000b7
)

func seccompStmt(code uint16, k uint32) syscall.SockFilter {
	return *syscall.LsfStmt(int(code), int(k))
}

func seccompJump(code uint16, k uint32, jt, jf uint8) syscall.SockFilter {
	return *syscall.LsfJump(int(code), int(k), int(jt), int(jf))
}

func seccompArch() (uint32, error) {
	switch runtime.GOARCH {
	case "amd64":
		return auditArchX86_64, nil
	case "arm64":
		return auditArchAArch64, nil
	default:
		return 0, fmt.Errorf("unsupported seccomp arch: %s", runtime.GOARCH)
	}
}

func defaultSeccompSyscalls() []uintptr {
	return []uintptr{
		unix.SYS_MOUNT,
		unix.SYS_UMOUNT2,
		unix.SYS_PIVOT_ROOT,
		unix.SYS_OPEN_TREE,
		unix.SYS_MOVE_MOUNT,
		unix.SYS_FSOPEN,
		unix.SYS_FSCONFIG,
		unix.SYS_FSMOUNT,
		unix.SYS_MOUNT_SETATTR,
		unix.SYS_UNSHARE,
		unix.SYS_SETNS,
		unix.SYS_CLONE3,
		unix.SYS_BPF,
		unix.SYS_PERF_EVENT_OPEN,
		unix.SYS_ADD_KEY,
		unix.SYS_REQUEST_KEY,
		unix.SYS_KEYCTL,
		unix.SYS_KEXEC_LOAD,
		unix.SYS_OPEN_BY_HANDLE_AT,
		unix.SYS_INIT_MODULE,
		unix.SYS_FINIT_MODULE,
		unix.SYS_DELETE_MODULE,
		unix.SYS_PTRACE,
		unix.SYS_PROCESS_VM_READV,
		unix.SYS_PROCESS_VM_WRITEV,
		unix.SYS_REBOOT,
		unix.SYS_SWAPON,
		unix.SYS_SWAPOFF,
		unix.SYS_SYSLOG,
	}
}

func installDefaultSeccomp() error {
	arch, err := seccompArch()
	if err != nil {
		return err
	}
	filters := []syscall.SockFilter{
		seccompStmt(unix.BPF_LD|unix.BPF_W|unix.BPF_ABS, seccompDataOffsetArch),
		seccompJump(unix.BPF_JMP|unix.BPF_JEQ|unix.BPF_K, arch, 1, 0),
		seccompStmt(unix.BPF_RET|unix.BPF_K, seccompRetKillProcess),
		seccompStmt(unix.BPF_LD|unix.BPF_W|unix.BPF_ABS, seccompDataOffsetNR),
	}
	for _, nr := range defaultSeccompSyscalls() {
		filters = append(filters,
			seccompJump(unix.BPF_JMP|unix.BPF_JEQ|unix.BPF_K, uint32(nr), 0, 1),
			seccompStmt(unix.BPF_RET|unix.BPF_K, seccompRetErrno|uint32(syscall.EPERM)),
		)
	}
	filters = append(filters, seccompStmt(unix.BPF_RET|unix.BPF_K, seccompRetAllow))
	prog := syscall.SockFprog{
		Len:    uint16(len(filters)),
		Filter: &filters[0],
	}
	_, _, errno := syscall.RawSyscall6(syscall.SYS_PRCTL,
		uintptr(syscall.PR_SET_SECCOMP),
		uintptr(seccompModeFilter),
		uintptr(unsafe.Pointer(&prog)),
		0, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func apparmorStackOnExec(profile string) error {
	if profile == "" {
		return fmt.Errorf("apparmor profile is required")
	}
	return os.WriteFile("/proc/self/attr/exec", []byte("stack "+profile), 0)
}

func init() { RegisterInstaller(installSyscallNS) }

func installSyscallNS() {
	// syscall/clone — (syscall/clone flags) → pid
	// Creates a new process with the given namespace flags.
	// flags is a bitwise OR of CLONE_* constants.
	cloneFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/clone expects 1 arg (flags)")
		}
		flags, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/clone expected Int flags")
		}
		pid, _, errno := syscall.RawSyscall(syscall.SYS_CLONE, uintptr(flags)|uintptr(syscall.SIGCHLD), 0, 0)
		if errno != 0 {
			return vm.NIL, fmt.Errorf("clone: %v", errno)
		}
		return vm.MakeInt(int(pid)), nil
	})

	// syscall/unshare — (syscall/unshare flags)
	unshareFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/unshare expects 1 arg (flags)")
		}
		flags, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/unshare expected Int flags")
		}
		if err := syscall.Unshare(int(flags)); err != nil {
			return vm.NIL, fmt.Errorf("unshare: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/mount — (syscall/mount source target fstype flags data)
	mountFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 5 {
			return vm.NIL, fmt.Errorf("syscall/mount expects 5 args (source target fstype flags data)")
		}
		source, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mount expected String source")
		}
		target, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mount expected String target")
		}
		fstype, ok := vs[2].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mount expected String fstype")
		}
		flags, ok := vs[3].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mount expected Int flags")
		}
		data, ok := vs[4].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mount expected String data")
		}
		if err := syscall.Mount(string(source), string(target), string(fstype), uintptr(flags), string(data)); err != nil {
			return vm.NIL, fmt.Errorf("mount: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/umount — (syscall/umount target flags)
	umountFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/umount expects 2 args (target flags)")
		}
		target, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/umount expected String target")
		}
		flags, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/umount expected Int flags")
		}
		if err := syscall.Unmount(string(target), int(flags)); err != nil {
			return vm.NIL, fmt.Errorf("umount: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/pivot-root — (syscall/pivot-root new-root put-old)
	pivotRootFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/pivot-root expects 2 args (new-root put-old)")
		}
		newRoot, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/pivot-root expected String new-root")
		}
		putOld, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/pivot-root expected String put-old")
		}
		if err := syscall.PivotRoot(string(newRoot), string(putOld)); err != nil {
			return vm.NIL, fmt.Errorf("pivot-root: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/chroot — (syscall/chroot path)
	chrootFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/chroot expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/chroot expected String path")
		}
		if err := syscall.Chroot(string(path)); err != nil {
			return vm.NIL, fmt.Errorf("chroot: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/chdir — (syscall/chdir path)
	chdirFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/chdir expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/chdir expected String path")
		}
		if err := syscall.Chdir(string(path)); err != nil {
			return vm.NIL, fmt.Errorf("chdir: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/sethostname — (syscall/sethostname name)
	sethostnameFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/sethostname expects 1 arg")
		}
		name, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/sethostname expected String")
		}
		if err := syscall.Sethostname([]byte(string(name))); err != nil {
			return vm.NIL, fmt.Errorf("sethostname: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/exec — (syscall/exec path argv env)
	// Replaces the current process. argv is a vector of strings, env is a vector of "K=V" strings.
	execFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("syscall/exec expects 3 args (path argv env)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/exec expected String path")
		}
		argvSeq, ok := vs[1].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/exec expected Sequable argv")
		}
		envSeq, ok := vs[2].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/exec expected Sequable env")
		}
		argv := seqToStrings(argvSeq.Seq())
		env := seqToStrings(envSeq.Seq())
		if err := syscall.Exec(string(path), argv, env); err != nil {
			return vm.NIL, fmt.Errorf("exec: %v", err)
		}
		return vm.NIL, nil // unreachable on success
	})

	// syscall/getpid — (syscall/getpid)
	getpidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(os.Getpid()), nil
	})

	// syscall/getuid — (syscall/getuid)
	getuidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(os.Getuid()), nil
	})

	// syscall/getgid — (syscall/getgid)
	getgidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(os.Getgid()), nil
	})

	// syscall/setuid — (syscall/setuid uid)
	setuidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/setuid expects 1 arg")
		}
		uid, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/setuid expected Int")
		}
		if err := syscall.Setuid(int(uid)); err != nil {
			return vm.NIL, fmt.Errorf("setuid: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/setgid — (syscall/setgid gid)
	setgidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/setgid expects 1 arg")
		}
		gid, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/setgid expected Int")
		}
		if err := syscall.Setgid(int(gid)); err != nil {
			return vm.NIL, fmt.Errorf("setgid: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/mkdir — (syscall/mkdir path mode)
	mkdirFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/mkdir expects 2 args (path mode)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mkdir expected String path")
		}
		mode, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mkdir expected Int mode")
		}
		if err := syscall.Mkdir(string(path), uint32(mode)); err != nil {
			return vm.NIL, fmt.Errorf("mkdir: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/mknod — (syscall/mknod path mode major minor)
	// Creates device nodes inside container filesystems.
	mknodFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 4 {
			return vm.NIL, fmt.Errorf("syscall/mknod expects 4 args (path mode major minor)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mknod expected String path")
		}
		mode, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mknod expected Int mode")
		}
		major, ok := vs[2].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mknod expected Int major")
		}
		minor, ok := vs[3].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mknod expected Int minor")
		}
		dev := int(unix.Mkdev(uint32(major), uint32(minor)))
		if err := unix.Mknod(string(path), uint32(mode), dev); err != nil {
			return vm.NIL, fmt.Errorf("mknod: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/waitpid — (syscall/waitpid pid options) → {:pid p :status s}
	waitpidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/waitpid expects 2 args (pid options)")
		}
		pid, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/waitpid expected Int pid")
		}
		opts, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/waitpid expected Int options")
		}
		var ws syscall.WaitStatus
		rpid, err := syscall.Wait4(int(pid), &ws, int(opts), nil)
		if err != nil {
			return vm.NIL, fmt.Errorf("waitpid: %v", err)
		}
		sig := 0
		if ws.Signaled() {
			sig = int(ws.Signal())
		}
		return waitResultMapping.StructToRecord(WaitResult{
			Pid:    rpid,
			Status: ws.ExitStatus(),
			Signal: sig,
		}), nil
	})

	// syscall/write-file — (syscall/write-file path content mode)
	// Useful for writing to /sys/fs/cgroup, /proc, etc.
	writeFileFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("syscall/write-file expects 3 args (path content mode)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/write-file expected String path")
		}
		content, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/write-file expected String content")
		}
		mode, ok := vs[2].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/write-file expected Int mode")
		}
		if err := os.WriteFile(string(path), []byte(string(content)), os.FileMode(mode)); err != nil {
			return vm.NIL, fmt.Errorf("write-file: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/read-file — (syscall/read-file path)
	readFileFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/read-file expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/read-file expected String path")
		}
		data, err := os.ReadFile(string(path))
		if err != nil {
			return vm.NIL, fmt.Errorf("read-file: %v", err)
		}
		return vm.String(data), nil
	})

	// syscall/mkdir-p — (syscall/mkdir-p path mode)
	mkdirpFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/mkdir-p expects 2 args (path mode)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mkdir-p expected String path")
		}
		mode, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/mkdir-p expected Int mode")
		}
		if err := os.MkdirAll(string(path), os.FileMode(mode)); err != nil {
			return vm.NIL, fmt.Errorf("mkdir-p: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/rmdir — (syscall/rmdir path)
	rmdirFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/rmdir expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/rmdir expected String path")
		}
		if err := syscall.Rmdir(string(path)); err != nil {
			return vm.NIL, fmt.Errorf("rmdir: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/rm-rf — (syscall/rm-rf path)
	rmrfFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/rm-rf expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/rm-rf expected String path")
		}
		if err := os.RemoveAll(string(path)); err != nil {
			return vm.NIL, fmt.Errorf("rm-rf: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/rm — (syscall/rm path)
	rmFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/rm expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/rm expected String path")
		}
		if err := os.Remove(string(path)); err != nil {
			return vm.NIL, fmt.Errorf("rm: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/symlink — (syscall/symlink target linkpath)
	symlinkFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/symlink expects 2 args (target linkpath)")
		}
		target, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/symlink expected String target")
		}
		linkpath, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/symlink expected String linkpath")
		}
		if err := os.Symlink(string(target), string(linkpath)); err != nil {
			return vm.NIL, fmt.Errorf("symlink: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/chmod — (syscall/chmod path mode)
	chmodFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/chmod expects 2 args (path mode)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/chmod expected String path")
		}
		mode, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/chmod expected Int mode")
		}
		if err := os.Chmod(string(path), os.FileMode(mode)); err != nil {
			return vm.NIL, fmt.Errorf("chmod: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/uname — (syscall/uname) → {:sysname "Linux" :machine "aarch64" :release "6.x" :nodename "host"}
	unameFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		var buf syscall.Utsname
		if err := syscall.Uname(&buf); err != nil {
			return vm.NIL, fmt.Errorf("uname: %v", err)
		}
		return unameResultMapping.StructToRecord(UnameResult{
			Sysname:  charsToString(buf.Sysname[:]),
			Machine:  charsToString(buf.Machine[:]),
			Release:  charsToString(buf.Release[:]),
			Nodename: charsToString(buf.Nodename[:]),
		}), nil
	})

	// syscall/pipe — (syscall/pipe) → [read-handle write-handle]
	// Creates an anonymous pipe. Both ends are IOHandles (*os.File-backed).
	pipeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("syscall/pipe expects 0 args")
		}
		r, w, err := os.Pipe()
		if err != nil {
			return vm.NIL, fmt.Errorf("pipe: %v", err)
		}
		return vm.ArrayVector([]vm.Value{
			vm.NewBoxed(NewIOHandle(r)),
			vm.NewBoxed(NewIOHandle(w)),
		}), nil
	})

	// syscall/signal-notify — (syscall/signal-notify ch sig...)
	// Deliver received signals as Ints onto ch. Starts a Go goroutine that
	// forwards from os/signal's channel until the target ch is closed.
	signalNotifyFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 {
			return vm.NIL, fmt.Errorf("syscall/signal-notify expects at least 2 args (ch sig...)")
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/signal-notify expected Chan as first arg")
		}
		sigs := make([]os.Signal, 0, len(vs)-1)
		for _, v := range vs[1:] {
			n, ok := v.(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("syscall/signal-notify expected Int signal")
			}
			sigs = append(sigs, syscall.Signal(int(n)))
		}
		goCh := make(chan os.Signal, 8)
		signal.Notify(goCh, sigs...)
		go func() {
			defer func() {
				// Swallow panic if ch was closed while we're sending.
				_ = recover()
				signal.Stop(goCh)
			}()
			for s := range goCh {
				ch <- vm.MakeInt(int(s.(syscall.Signal)))
			}
		}()
		return vm.NIL, nil
	})

	// syscall/kill — (syscall/kill pid sig)
	killFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("syscall/kill expects 2 args (pid sig)")
		}
		pid, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/kill expected Int pid")
		}
		sig, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/kill expected Int sig")
		}
		if err := syscall.Kill(int(pid), syscall.Signal(int(sig))); err != nil {
			return vm.NIL, fmt.Errorf("kill: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/prctl — (syscall/prctl option arg2 arg3 arg4 arg5)
	// Thin prctl(2) wrapper for child-process setup like no_new_privs.
	prctlFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 5 {
			return vm.NIL, fmt.Errorf("syscall/prctl expects 5 args (option arg2 arg3 arg4 arg5)")
		}
		args := make([]uintptr, 5)
		for i, v := range vs {
			n, ok := v.(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("syscall/prctl expected Int at arg %d", i+1)
			}
			args[i] = uintptr(n)
		}
		_, _, errno := syscall.Syscall6(syscall.SYS_PRCTL,
			args[0], args[1], args[2], args[3], args[4], 0)
		if errno != 0 {
			return vm.NIL, fmt.Errorf("prctl: %v", errno)
		}
		return vm.NIL, nil
	})

	// syscall/capset — (syscall/capset effective permitted inheritable)
	// Sets VFS capabilities using Linux capability version 3 masks.
	capsetFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("syscall/capset expects 3 args (effective permitted inheritable)")
		}
		masks := make([]uint64, 3)
		for i, v := range vs {
			n, ok := v.(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("syscall/capset expected Int at arg %d", i+1)
			}
			masks[i] = uint64(n)
		}
		type capUserHeader struct {
			Version uint32
			Pid     int32
		}
		type capUserData struct {
			Effective   uint32
			Permitted   uint32
			Inheritable uint32
		}
		hdr := capUserHeader{Version: 0x20080522}
		data := [2]capUserData{
			{
				Effective:   uint32(masks[0]),
				Permitted:   uint32(masks[1]),
				Inheritable: uint32(masks[2]),
			},
			{
				Effective:   uint32(masks[0] >> 32),
				Permitted:   uint32(masks[1] >> 32),
				Inheritable: uint32(masks[2] >> 32),
			},
		}
		_, _, errno := syscall.RawSyscall(syscall.SYS_CAPSET,
			uintptr(unsafe.Pointer(&hdr)),
			uintptr(unsafe.Pointer(&data[0])),
			0)
		if errno != 0 {
			return vm.NIL, fmt.Errorf("capset: %v", errno)
		}
		return vm.NIL, nil
	})

	// syscall/seccomp-default — (syscall/seccomp-default)
	// Installs lgcr's built-in default seccomp filter.
	seccompDefaultFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("syscall/seccomp-default expects 0 args")
		}
		if err := installDefaultSeccomp(); err != nil {
			return vm.NIL, fmt.Errorf("seccomp-default: %v", err)
		}
		return vm.NIL, nil
	})

	// syscall/apparmor-stack-onexec — (syscall/apparmor-stack-onexec profile)
	// Schedules an AppArmor profile stack transition at the next exec boundary.
	apparmorStackOnExecFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("syscall/apparmor-stack-onexec expects 1 arg (profile)")
		}
		profile, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/apparmor-stack-onexec expected String profile")
		}
		if err := apparmorStackOnExec(string(profile)); err != nil {
			return vm.NIL, fmt.Errorf("apparmor-stack-onexec: %v", err)
		}
		return vm.NIL, nil
	})

	// stdioFile resolves a spawn stdio slot value to an *os.File.
	// nil → opens /dev/null with the given flag.
	// IOHandle → returns underlying File.
	stdioFile := func(v vm.Value, flag int) (*os.File, bool, error) {
		if v == vm.NIL || v == nil {
			f, err := os.OpenFile(os.DevNull, flag, 0)
			if err != nil {
				return nil, false, err
			}
			return f, true, nil // caller must close
		}
		h, err := getIOHandle(v)
		if err != nil {
			return nil, false, err
		}
		f := h.File()
		if f == nil {
			return nil, false, fmt.Errorf("IOHandle is not file-backed; spawn stdio requires a real file descriptor")
		}
		return f, false, nil
	}

	// syscall/spawn-async — (syscall/spawn-async path argv env cloneflags stdin stdout stderr [opts])
	// Non-blocking: returns {:pid p} immediately. Caller waits via syscall/waitpid.
	// Each stdio slot is nil (→ /dev/null) or an IOHandle. The child gets a dup of
	// the underlying fd; the parent retains its handle and may close it after spawn.
	// Optional 8th arg: map of options.
	//   {:setctty? true} — make stdin the child's controlling terminal (pty slave).
	spawnAsyncFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 7 || len(vs) > 8 {
			return vm.NIL, fmt.Errorf("syscall/spawn-async expects 7-8 args (path argv env cloneflags stdin stdout stderr [opts])")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn-async expected String path")
		}
		argvSeq, ok := vs[1].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn-async expected Sequable argv")
		}
		envSeq, ok := vs[2].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn-async expected Sequable env")
		}
		flags, ok := vs[3].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn-async expected Int cloneflags")
		}
		argv := seqToStrings(argvSeq.Seq())
		env := seqToStrings(envSeq.Seq())

		stdinF, stdinOwn, err := stdioFile(vs[4], os.O_RDONLY)
		if err != nil {
			return vm.NIL, fmt.Errorf("spawn-async stdin: %v", err)
		}
		if stdinOwn {
			defer stdinF.Close()
		}
		stdoutF, stdoutOwn, err := stdioFile(vs[5], os.O_WRONLY)
		if err != nil {
			return vm.NIL, fmt.Errorf("spawn-async stdout: %v", err)
		}
		if stdoutOwn {
			defer stdoutF.Close()
		}
		stderrF, stderrOwn, err := stdioFile(vs[6], os.O_WRONLY)
		if err != nil {
			return vm.NIL, fmt.Errorf("spawn-async stderr: %v", err)
		}
		if stderrOwn {
			defer stderrF.Close()
		}

		sysAttr := &syscall.SysProcAttr{
			Cloneflags: uintptr(flags),
			// Give the child its own session so it survives the parent's
			// shell/sudo tearing down the process group.
			Setsid: true,
		}
		var dir string
		if len(vs) == 8 && vs[7] != vm.NIL {
			optsMap, ok := vs[7].(*vm.PersistentMap)
			if !ok {
				return vm.NIL, fmt.Errorf("spawn-async opts must be a map")
			}
			if v := optsMap.ValueAt(vm.Keyword("setctty?")); v != nil && v != vm.NIL && v != vm.FALSE {
				// stdin (fd 0 in the child after the Files[] dup) becomes the
				// controlling terminal. Requires Setsid (already set).
				sysAttr.Setctty = true
				sysAttr.Ctty = 0
			}
			// :uid / :gid — drop privileges between fork and exec via
			// SysProcAttr.Credential. Both must be set together.
			uidV := optsMap.ValueAt(vm.Keyword("uid"))
			gidV := optsMap.ValueAt(vm.Keyword("gid"))
			if uidV != nil && uidV != vm.NIL {
				u, ok := uidV.(vm.Int)
				if !ok {
					return vm.NIL, fmt.Errorf("spawn-async :uid must be Int")
				}
				cred := &syscall.Credential{Uid: uint32(u), NoSetGroups: true}
				if gidV != nil && gidV != vm.NIL {
					g, ok := gidV.(vm.Int)
					if !ok {
						return vm.NIL, fmt.Errorf("spawn-async :gid must be Int")
					}
					cred.Gid = uint32(g)
				} else {
					cred.Gid = uint32(u)
				}
				sysAttr.Credential = cred
			}
			// :dir — chdir in the child before exec (Go's os.StartProcess Dir).
			if d := optsMap.ValueAt(vm.Keyword("dir")); d != nil && d != vm.NIL {
				ds, ok := d.(vm.String)
				if !ok {
					return vm.NIL, fmt.Errorf("spawn-async :dir must be String")
				}
				dir = string(ds)
			}
		}
		proc, err := os.StartProcess(string(path), argv, &os.ProcAttr{
			Dir:   dir,
			Env:   env,
			Files: []*os.File{stdinF, stdoutF, stderrF},
			Sys:   sysAttr,
		})
		if err != nil {
			return vm.NIL, fmt.Errorf("spawn-async: %v", err)
		}
		// Release the *os.Process — waitpid will reap it via the kernel.
		pid := proc.Pid
		_ = proc.Release()
		return spawnResultMapping.StructToRecord(SpawnResult{Pid: pid}), nil
	})

	// syscall/spawn — (syscall/spawn path argv env cloneflags)
	// Fork+exec a child process with namespace cloneflags.
	// Returns {:pid p :exit e :out "..." :err "..."} after the child exits.
	spawnFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 4 {
			return vm.NIL, fmt.Errorf("syscall/spawn expects 4 args (path argv env cloneflags)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn expected String path")
		}
		argvSeq, ok := vs[1].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn expected Sequable argv")
		}
		envSeq, ok := vs[2].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn expected Sequable env")
		}
		flags, ok := vs[3].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("syscall/spawn expected Int cloneflags")
		}
		argv := seqToStrings(argvSeq.Seq())
		env := seqToStrings(envSeq.Seq())

		cmd := exec.Command(string(path), argv[1:]...)
		cmd.Env = env
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: uintptr(flags),
		}
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return vm.NIL, fmt.Errorf("spawn: %v", err)
			}
		}
		return spawnResultMapping.StructToRecord(SpawnResult{
			Pid:  cmd.ProcessState.Pid(),
			Exit: exitCode,
			Out:  stdout.String(),
			Err:  stderr.String(),
		}), nil
	})

	ns := vm.NewNamespace("syscall")

	// Intentional shadows of clojure.core names — suppress warn-on-shadow.
	ns.Exclude("mkdir")

	// namespace creation
	ns.Def("clone", cloneFn)
	ns.Def("unshare", unshareFn)

	// filesystem
	ns.Def("mount", mountFn)
	ns.Def("umount", umountFn)
	ns.Def("pivot-root", pivotRootFn)
	ns.Def("chroot", chrootFn)
	ns.Def("chdir", chdirFn)
	ns.Def("mkdir", mkdirFn)
	ns.Def("mknod", mknodFn)
	ns.Def("mkdir-p", mkdirpFn)
	ns.Def("rmdir", rmdirFn)
	ns.Def("rm-rf", rmrfFn)
	ns.Def("rm", rmFn)
	ns.Def("symlink", symlinkFn)
	ns.Def("chmod", chmodFn)

	// hostname
	ns.Def("sethostname", sethostnameFn)

	// process
	ns.Def("exec", execFn)
	ns.Def("spawn", spawnFn)
	ns.Def("spawn-async", spawnAsyncFn)
	ns.Def("pipe", pipeFn)
	ns.Def("kill", killFn)
	ns.Def("prctl", prctlFn)
	ns.Def("capset", capsetFn)
	ns.Def("seccomp-default", seccompDefaultFn)
	ns.Def("apparmor-stack-onexec", apparmorStackOnExecFn)
	ns.Def("signal-notify", signalNotifyFn)
	ns.Def("getpid", getpidFn)
	ns.Def("getuid", getuidFn)
	ns.Def("getgid", getgidFn)
	ns.Def("setuid", setuidFn)
	ns.Def("setgid", setgidFn)
	ns.Def("waitpid", waitpidFn)

	// system info
	ns.Def("uname", unameFn)

	// file I/O (for /proc, /sys/fs/cgroup, etc.)
	ns.Def("write-file", writeFileFn)
	ns.Def("read-file", readFileFn)

	// clone flags
	ns.Def("CLONE_NEWNS", vm.MakeInt(0x00020000))
	ns.Def("CLONE_NEWUTS", vm.MakeInt(0x04000000))
	ns.Def("CLONE_NEWIPC", vm.MakeInt(0x08000000))
	ns.Def("CLONE_NEWPID", vm.MakeInt(0x20000000))
	ns.Def("CLONE_NEWNET", vm.MakeInt(0x40000000))
	ns.Def("CLONE_NEWUSER", vm.MakeInt(0x10000000))

	// mount flags
	ns.Def("MS_BIND", vm.MakeInt(4096))
	ns.Def("MS_REC", vm.MakeInt(16384))
	ns.Def("MS_PRIVATE", vm.MakeInt(1<<18))
	ns.Def("MS_RDONLY", vm.MakeInt(1))
	ns.Def("MS_NOSUID", vm.MakeInt(2))
	ns.Def("MS_NODEV", vm.MakeInt(4))
	ns.Def("MS_NOEXEC", vm.MakeInt(8))

	ns.Def("S_IFCHR", vm.MakeInt(0o020000))

	// signals
	ns.Def("SIGHUP", vm.MakeInt(1))
	ns.Def("SIGINT", vm.MakeInt(2))
	ns.Def("SIGQUIT", vm.MakeInt(3))
	ns.Def("SIGKILL", vm.MakeInt(9))
	ns.Def("SIGTERM", vm.MakeInt(15))
	ns.Def("SIGCHLD", vm.MakeInt(17))
	ns.Def("SIGWINCH", vm.MakeInt(28))

	// waitpid options
	ns.Def("WNOHANG", vm.MakeInt(1))

	// prctl options
	ns.Def("PR_CAPBSET_DROP", vm.MakeInt(24))
	ns.Def("PR_SET_NO_NEW_PRIVS", vm.MakeInt(38))

	RegisterNS(ns)
}

// charsToString converts a null-terminated char array (from Utsname) to a string.
// Utsname fields are []int8 on some Linux arches (amd64) and []uint8 on others (arm).
func charsToString[T int8 | uint8](ca []T) string {
	buf := make([]byte, 0, len(ca))
	for _, c := range ca {
		if c == 0 {
			break
		}
		buf = append(buf, byte(c))
	}
	return string(buf)
}

// seqToStrings converts a Seq to a []string.
func seqToStrings(s vm.Seq) []string {
	var result []string
	for s != nil {
		v := s.First()
		if str, ok := v.(vm.String); ok {
			result = append(result, string(str))
		} else {
			result = append(result, v.String())
		}
		s = s.Next()
	}
	return result
}

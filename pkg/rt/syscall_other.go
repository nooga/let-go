//go:build !linux

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installSyscallNS) }

func installSyscallNS() {
	unsupported := func(name string) vm.Value {
		fn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
			return vm.NIL, fmt.Errorf("syscall/%s is only supported on Linux", name)
		})
		return fn
	}

	// syscall/getpid — works everywhere
	getpidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(os.Getpid()), nil
	})

	// syscall/getuid — works everywhere
	getuidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(os.Getuid()), nil
	})

	// syscall/getgid — works everywhere
	getgidFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(os.Getgid()), nil
	})

	// syscall/read-file — works everywhere
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

	// syscall/write-file — works everywhere
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

	// syscall/mkdir-p — works everywhere
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

	// syscall/rm-rf — works everywhere
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

	// syscall/rm — works everywhere
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

	// syscall/symlink — works everywhere
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

	// syscall/chmod — works everywhere
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

	ns := vm.NewNamespace("syscall")

	// Intentional shadows of clojure.core names — suppress warn-on-shadow.
	ns.Exclude("mkdir")

	// Linux-only stubs
	ns.Def("clone", unsupported("clone"))
	ns.Def("unshare", unsupported("unshare"))
	ns.Def("mount", unsupported("mount"))
	ns.Def("umount", unsupported("umount"))
	ns.Def("pivot-root", unsupported("pivot-root"))
	ns.Def("chroot", unsupported("chroot"))
	ns.Def("chdir", unsupported("chdir"))
	ns.Def("mkdir", unsupported("mkdir"))
	ns.Def("mknod", unsupported("mknod"))
	ns.Def("rmdir", unsupported("rmdir"))
	ns.Def("sethostname", unsupported("sethostname"))
	ns.Def("exec", unsupported("exec"))
	ns.Def("spawn", unsupported("spawn"))
	ns.Def("spawn-async", unsupported("spawn-async"))
	ns.Def("pipe", unsupported("pipe"))
	ns.Def("kill", unsupported("kill"))
	ns.Def("prctl", unsupported("prctl"))
	ns.Def("capset", unsupported("capset"))
	ns.Def("seccomp-default", unsupported("seccomp-default"))
	ns.Def("apparmor-stack-onexec", unsupported("apparmor-stack-onexec"))
	ns.Def("signal-notify", unsupported("signal-notify"))
	ns.Def("uname", unsupported("uname"))
	ns.Def("setuid", unsupported("setuid"))
	ns.Def("setgid", unsupported("setgid"))
	ns.Def("waitpid", unsupported("waitpid"))

	// Cross-platform
	ns.Def("getpid", getpidFn)
	ns.Def("getuid", getuidFn)
	ns.Def("getgid", getgidFn)
	ns.Def("read-file", readFileFn)
	ns.Def("write-file", writeFileFn)
	ns.Def("mkdir-p", mkdirpFn)
	ns.Def("rm-rf", rmrfFn)
	ns.Def("rm", rmFn)
	ns.Def("symlink", symlinkFn)
	ns.Def("chmod", chmodFn)

	// Constants (same values everywhere so code compiles/loads on any platform)
	ns.Def("CLONE_NEWNS", vm.MakeInt(0x00020000))
	ns.Def("CLONE_NEWUTS", vm.MakeInt(0x04000000))
	ns.Def("CLONE_NEWIPC", vm.MakeInt(0x08000000))
	ns.Def("CLONE_NEWPID", vm.MakeInt(0x20000000))
	ns.Def("CLONE_NEWNET", vm.MakeInt(0x40000000))
	ns.Def("CLONE_NEWUSER", vm.MakeInt(0x10000000))

	ns.Def("MS_BIND", vm.MakeInt(4096))
	ns.Def("MS_REC", vm.MakeInt(16384))
	ns.Def("MS_PRIVATE", vm.MakeInt(1<<18))
	ns.Def("MS_RDONLY", vm.MakeInt(1))
	ns.Def("MS_NOSUID", vm.MakeInt(2))
	ns.Def("MS_NODEV", vm.MakeInt(4))
	ns.Def("MS_NOEXEC", vm.MakeInt(8))

	ns.Def("S_IFCHR", vm.MakeInt(0o020000))

	ns.Def("WNOHANG", vm.MakeInt(1))
	ns.Def("PR_CAPBSET_DROP", vm.MakeInt(24))
	ns.Def("PR_SET_NO_NEW_PRIVS", vm.MakeInt(38))

	// signals (Linux values — constants load everywhere, call sites error)
	ns.Def("SIGHUP", vm.MakeInt(1))
	ns.Def("SIGINT", vm.MakeInt(2))
	ns.Def("SIGQUIT", vm.MakeInt(3))
	ns.Def("SIGKILL", vm.MakeInt(9))
	ns.Def("SIGTERM", vm.MakeInt(15))
	ns.Def("SIGCHLD", vm.MakeInt(17))
	ns.Def("SIGWINCH", vm.MakeInt(28))

	RegisterNS(ns)
}

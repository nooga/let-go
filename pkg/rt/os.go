//go:build !tinygo

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installOsNS) }

// nolint
func installOsNS() {
	getenv, err := vm.NativeFnType.Box(os.Getenv)
	execf, err := vm.NativeFnType.Box(exec.Command)
	tempDir, err := vm.NativeFnType.Box(os.TempDir)
	args, err := vm.ToLetGo(os.Args)
	withStdin, err := vm.NativeFnType.Wrap(func(v []vm.Value) (vm.Value, error) {
		var cmd = v[0].Unbox().(*exec.Cmd)
		s := string(v[1].(vm.String))
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return vm.NIL, err
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, s)
		}()
		return v[0], nil
	})

	// os/exit — (os/exit code)
	exitf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/exit expects 1 arg")
		}
		code, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("os/exit expected Int")
		}
		RunExitHooks() // flush profiling etc.; os.Exit skips deferred funcs
		os.Exit(int(code))
		return vm.NIL, nil
	})

	// os/cwd — (os/cwd)
	cwd, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		d, err := os.Getwd()
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(d), nil
	})

	// os/setenv — (os/setenv key val)
	setenv, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("os/setenv expects 2 args")
		}
		k, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/setenv expected String key")
		}
		v, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/setenv expected String value")
		}
		return vm.NIL, os.Setenv(string(k), string(v))
	})

	// os/ls — (os/ls path)
	ls, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/ls expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/ls expected String path")
		}
		entries, err := os.ReadDir(string(path))
		if err != nil {
			return vm.NIL, err
		}
		result := make([]vm.Value, len(entries))
		for i, e := range entries {
			result[i] = vm.String(e.Name())
		}
		return vm.NewArrayVector(result), nil
	})

	// os/stat — (os/stat path)
	stat, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/stat expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/stat expected String path")
		}
		info, err := os.Stat(string(path))
		if err != nil {
			if os.IsNotExist(err) {
				return vm.NIL, nil
			}
			return vm.NIL, err
		}
		return fileStatMapping.StructToRecord(FileStat{
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime().String(),
		}), nil
	})

	// os/sh — (os/sh cmd & args) → {:exit 0 :out "..." :err "..."}
	sh, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("os/sh expects at least 1 arg")
		}
		cmdName, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/sh expected String command")
		}
		args := make([]string, len(vs)-1)
		for i := 1; i < len(vs); i++ {
			if s, ok := vs[i].(vm.String); ok {
				args[i-1] = string(s)
			} else {
				args[i-1] = vs[i].String()
			}
		}
		cmd := exec.Command(string(cmdName), args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return vm.NIL, err
			}
		}
		return shellResultMapping.StructToRecord(ShellResult{
			Exit: exitCode,
			Out:  stdout.String(),
			Err:  stderr.String(),
		}), nil
	})

	// os/exec* — (os/exec* cmd & args) → exit code (Int). Unlike os/sh (which
	// buffers and returns {:out :err :exit}), the child STREAMS: its stdout and
	// stderr are wired to the current dynamic bindings of *out* and *err* — the
	// same streams println / (binding [*out* ...] ...) / with-out-str drive — so
	// a caller that redirects or captures those vars actually sees the child's
	// output. (Previously the child was pinned to the raw os.Stdout/os.Stderr,
	// which escaped every rebinding: a build harness capturing only *out* lost
	// the child's stderr entirely.) Falls back to the process streams when the
	// vars aren't installed (early boot). Stdin stays os.Stdin so the child can
	// stay interactive (e.g. launching a REPL). Returns the exit code.
	execStar := vm.NewCtxNativeFn("exec*", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("os/exec* expects at least 1 arg")
		}
		cmdName, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/exec* expected String command")
		}
		cmdArgs := make([]string, len(vs)-1)
		for i := 1; i < len(vs); i++ {
			if s, ok := vs[i].(vm.String); ok {
				cmdArgs[i-1] = string(s)
			} else {
				cmdArgs[i-1] = vs[i].String()
			}
		}
		cmd := exec.Command(string(cmdName), cmdArgs...)
		cmd.Stdin = os.Stdin
		// Wire child stdout/stderr to the current *out*/*err* writers, mirroring
		// how println resolves *out* (resolveIOHandleVar respects binding).
		outH := resolveIOHandleVar(ec, "*out*")
		errH := resolveIOHandleVar(ec, "*err*")
		if outH != nil && outH.Writer() != nil {
			cmd.Stdout = outH.Writer()
		} else {
			cmd.Stdout = os.Stdout
		}
		if errH != nil && errH.Writer() != nil {
			cmd.Stderr = errH.Writer()
		} else {
			cmd.Stderr = os.Stderr
		}
		exitCode := 0
		runErr := cmd.Run()
		// Flush buffered handles (bufio-backed *out*/*err*) so the child's output
		// is visible before we hand control back to the caller.
		if outH != nil {
			_ = outH.Sync()
		}
		if errH != nil {
			_ = errH.Sync()
		}
		if runErr != nil {
			if exitErr, ok := runErr.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return vm.NIL, runErr
			}
		}
		return vm.Int(exitCode), nil
	})

	if err != nil {
		panic(fmt.Sprintf("os NS init failed: %e", err))
	}

	ns := vm.NewNamespace("os")

	ns.Def("getenv", getenv)
	ns.Def("exec", execf)
	ns.Def("with-stdin", withStdin)
	ns.Def("temp-dir", tempDir)
	ns.Def("args", args)
	ns.Def("exit", exitf)
	ns.Def("cwd", cwd)
	ns.Def("setenv", setenv)
	ns.Def("ls", ls)
	ns.Def("stat", stat)
	ns.Def("sh", sh)
	ns.Def("exec*", execStar)

	// os/free-port — (os/free-port) → an OS-assigned free TCP port (Int).
	// Binds 127.0.0.1:0, reads the assigned port, and releases the listener.
	// Check-then-use: the port could in principle be taken between this call
	// and the caller's own bind, so treat it as a strong hint rather than a
	// reservation (kernels avoid immediately reissuing a just-released port).
	ns.Def("free-port", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("os/free-port expects no args")
		}
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return vm.NIL, err
		}
		port := l.Addr().(*net.TCPAddr).Port
		if err := l.Close(); err != nil {
			return vm.NIL, err
		}
		return vm.Int(port), nil
	}))

	// os/unzip — (os/unzip zip-path dest-dir) → dest-dir
	// Extracts a zip archive into dest-dir, creating it if missing and
	// overwriting existing files. Entries that would land outside dest-dir
	// are refused (see unzipEntryTarget); symlink entries are skipped.
	ns.Def("unzip", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("os/unzip expects 2 args")
		}
		src, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/unzip expected String path")
		}
		dest, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/unzip expected String destination")
		}
		if err := unzipTo(string(src), string(dest)); err != nil {
			return vm.NIL, err
		}
		return vs[1], nil
	}))

	// os/rename — (os/rename old-path new-path) → new-path
	//
	// Renames old-path to new-path. On Unix, a rename within one filesystem
	// is atomic: a concurrent reader sees either the old state or the new one,
	// never a half-written file. Go does not guarantee that property on
	// non-Unix hosts. Writing to a temporary name in the destination directory
	// and renaming into place is the usual Unix way to publish a file safely.
	//
	// A rename across filesystems fails rather than falling back to
	// copy-then-delete. The fallback is what callers reach for this instead
	// of, so silently substituting it would remove the only property that
	// distinguishes it from spit.
	ns.Def("rename", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("os/rename expects 2 args")
		}
		from, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/rename expected String path")
		}
		to, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/rename expected String destination")
		}
		if err := os.Rename(string(from), string(to)); err != nil {
			return vm.NIL, err
		}
		return vs[1], nil
	}))

	// os/delete-tree — (os/delete-tree path) → nil
	//
	// Removes path and everything beneath it: the recursive form of
	// delete-file, which removes a single entry and fails on a non-empty
	// directory.
	//
	// syscall/rm-rf is the same RemoveAll call and predates this. It lives
	// in a namespace built for container setup, beside clone, pivot-root
	// and seccomp, which is not where a caller doing ordinary filesystem
	// work looks. This is the os-namespace spelling, and it is the one that
	// carries the guards below.
	//
	// Three behaviours worth knowing, all inherited from RemoveAll:
	//
	//   - Removing something already absent succeeds. The post-state the
	//     caller asked for is the one that holds. delete-file errors here.
	//   - Symlinks are unlinked, never followed, so a link inside the tree
	//     pointing outside it does not take the target down with it.
	//   - By the same rule, a path that is *itself* a symlink to a
	//     directory loses only the link; the directory it names keeps its
	//     contents. Callers wanting the target gone should canonicalize
	//     first.
	//
	// An empty path is refused: RemoveAll treats "" as a silent no-op,
	// which hides the unset variable that produced it. The filesystem root
	// is refused for the same reason and a worse outcome — (str nil) is ""
	// in lg, so a caller building "$root/$name" with root unset produces
	// "/name", one level below the root it never meant to touch. RemoveAll
	// already rejects "." itself.
	ns.Def("delete-tree", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/delete-tree expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/delete-tree expected String path")
		}
		if path == "" {
			return vm.NIL, fmt.Errorf("os/delete-tree expected a non-empty path")
		}
		// Dir(p) == p exactly at a volume root: "/" on unix, "C:\" and "\"
		// on Windows. Cheaper and more portable than matching separators.
		if cleaned := filepath.Clean(string(path)); filepath.Dir(cleaned) == cleaned {
			return vm.NIL, fmt.Errorf("os/delete-tree refused the filesystem root %q", cleaned)
		}
		if err := os.RemoveAll(string(path)); err != nil {
			return vm.NIL, err
		}
		return vm.NIL, nil
	}))

	// os/absolute-path — (os/absolute-path path) → "/abs/path"
	//
	// Resolves path against the process working directory and cleans it.
	// Lexical with respect to the argument: the path itself is never looked
	// up, so it need not exist and any symlink in it stays a symlink. (The
	// process cwd is read, which is why this can still fail — os.Getwd
	// errors when the working directory has been removed.)
	//
	// An empty path is refused rather than quietly meaning the cwd, on the
	// same reasoning as delete-tree: it is what an unset variable looks
	// like, and returning a plausible answer hides it.
	ns.Def("absolute-path", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/absolute-path expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/absolute-path expected String path")
		}
		if path == "" {
			return vm.NIL, fmt.Errorf("os/absolute-path expected a non-empty path")
		}
		abs, err := filepath.Abs(string(path))
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(abs), nil
	}))

	// os/canonical-path — (os/canonical-path path) → "/real/path"
	//
	// The absolute path with every symlink resolved, so two names for one
	// file produce one string. That is what makes it the form to compare or
	// use as a key, and it is why this reads the filesystem where
	// absolute-path does not: a symlink can only be followed by looking.
	//
	// A path that does not exist is an error rather than a cleaned-up
	// guess. Callers wanting a name for a file they are about to create
	// want absolute-path.
	ns.Def("canonical-path", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/canonical-path expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/canonical-path expected String path")
		}
		if path == "" {
			return vm.NIL, fmt.Errorf("os/canonical-path expected a non-empty path")
		}
		// Resolve before absolutizing. filepath.Abs cleans lexically, which
		// collapses ".." against the *link's* parent instead of the
		// target's — so Abs-then-EvalSymlinks reports ENOENT for a path the
		// kernel opens fine. EvalSymlinks walks a relative path against the
		// cwd itself, so doing it first costs nothing.
		real, err := filepath.EvalSymlinks(string(path))
		if err != nil {
			return vm.NIL, err
		}
		abs, err := filepath.Abs(real)
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(abs), nil
	}))

	// os/os-name — (os/os-name) → "linux", "darwin", "windows", ...
	ns.Def("os-name", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.String(runtime.GOOS), nil
	}))

	// os/arch — (os/arch) → "amd64", "arm64", ...
	ns.Def("arch", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.String(runtime.GOARCH), nil
	}))

	// os/user-name — (os/user-name)
	ns.Def("user-name", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if _, name := currentUser(); name != "" {
			return vm.String(name), nil
		}
		return vm.String(os.Getenv("USER")), nil
	}))

	// os/hostname — (os/hostname)
	ns.Def("hostname", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		h, err := os.Hostname()
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(h), nil
	}))

	// os/file-separator, os/path-separator, os/line-separator
	ns.Def("file-separator", vm.String(string(os.PathSeparator)))
	ns.Def("path-separator", vm.String(string(os.PathListSeparator)))
	ns.Def("line-separator", vm.String(lineSeparator()))

	RegisterNS(ns)
}

func lineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// unzipTo extracts src into dest.
//
// Every write goes through an os.Root confined to dest, so an entry can
// neither escape lexically ("../evil.txt") nor through a symlink that already
// exists inside dest ("link/x" where dest/link points elsewhere). os.Root
// enforces containment while it walks each path component, which a
// check-then-write guard fundamentally cannot: another process sharing dest
// could swap a validated directory for a symlink in the window between the
// check and the write.
//
// Fidelity is traded for safety besides: symlink entries and other
// non-regular entries (devices, fifos, sockets) are skipped rather than
// recreated. Permissions follow the unzip(1) contract — as recorded in the
// entry, masked by the process umask.
func unzipTo(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	root, err := os.OpenRoot(dest)
	if err != nil {
		return err
	}
	defer root.Close()

	for _, f := range r.File {
		mode := f.Mode()
		if mode&os.ModeSymlink != 0 {
			continue
		}
		// Zip names are always slash-separated, whatever the host.
		name := filepath.Clean(filepath.FromSlash(f.Name))
		if name == "." {
			continue
		}
		if f.FileInfo().IsDir() {
			err = unzipDir(root, name, unzipDirPerm(f))
		} else if mode.IsRegular() {
			err = unzipFile(root, f, name)
		} else {
			continue
		}
		if err != nil {
			return fmt.Errorf("os/unzip: %s: %w", f.Name, err)
		}
	}
	return nil
}

// unzipEntryPerm reads the permissions a zip entry records for itself. Zip
// keeps unix permissions in the high 16 bits of the external attributes; a
// DOS/FAT-style writer leaves them empty and archive/zip synthesizes 0666,
// which under a permissive umask would land world-writable files on disk. So
// distinguish "recorded" from "synthesized" here and let callers pick their
// own default for the latter.
func unzipEntryPerm(f *zip.File) (os.FileMode, bool) {
	if unix := f.ExternalAttrs >> 16; unix != 0 {
		if perm := os.FileMode(unix).Perm(); perm != 0 {
			return perm, true
		}
	}
	return 0, false
}

// unzipFilePerm is the mode to create a file entry with — as recorded (this
// is what carries the executable bit), else 0644.
func unzipFilePerm(f *zip.File) os.FileMode {
	if perm, ok := unzipEntryPerm(f); ok {
		return perm
	}
	return 0o644
}

// unzipDirPerm is the mode to create a directory entry with. Owner rwx is
// forced on: a read-only directory entry listed before the files it contains
// would otherwise make the rest of the archive unextractable.
func unzipDirPerm(f *zip.File) os.FileMode {
	if perm, ok := unzipEntryPerm(f); ok {
		return perm | 0o700
	}
	return 0o755
}

// unzipParents creates name's parent directories inside root.
func unzipParents(root *os.Root, name string) error {
	parent := filepath.Dir(name)
	if parent == "." || parent == string(filepath.Separator) {
		return nil
	}
	return root.MkdirAll(parent, 0o755)
}

// unzipDir materialises an explicit directory entry. A directory created
// earlier as some file's implicit parent keeps the 0755 it got then — only a
// freshly created one carries the entry's recorded mode.
func unzipDir(root *os.Root, name string, perm os.FileMode) error {
	if err := unzipParents(root, name); err != nil {
		return err
	}
	err := root.Mkdir(name, perm)
	if err == nil || !errors.Is(err, os.ErrExist) {
		return err
	}
	// "Already exists" is only benign when what exists is itself a directory.
	// A plain file (or a symlink) sitting at a directory entry's path is a
	// genuine conflict — swallowing it would report a successful extraction
	// while leaving the entry's contents nowhere to go.
	info, statErr := root.Lstat(name)
	if statErr != nil {
		return statErr
	}
	if !info.IsDir() {
		return errors.New("exists and is not a directory")
	}
	return nil
}

// unzipFile writes one regular entry, creating parents as needed.
func unzipFile(root *os.Root, f *zip.File, name string) error {
	if err := unzipParents(root, name); err != nil {
		return err
	}
	// Open the entry BEFORE touching the destination: f.Open fails outright
	// on an unsupported compression method, and destroying a perfectly good
	// existing file on the way to an error nobody can recover from is worse
	// than not extracting at all.
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Clear whatever is already at the target. Removing rather than
	// overwriting in place matters twice over: O_CREATE applies its perm
	// argument only when it creates the file, so an in-place overwrite would
	// silently keep the old file's mode; and a symlink sitting there would
	// otherwise be followed (os.Root confines where it may point, but a link
	// to another path inside dest is legal and would still be written
	// through).
	if info, err := root.Lstat(name); err == nil && !info.IsDir() {
		if err := root.Remove(name); err != nil {
			return err
		}
	}

	out, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, unzipFilePerm(f))
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, rc); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func mustWrap(fn func([]vm.Value) (vm.Value, error)) vm.Value {
	v, err := vm.NativeFnType.Wrap(fn)
	if err != nil {
		panic(err)
	}
	return v
}

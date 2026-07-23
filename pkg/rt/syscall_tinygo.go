//go:build tinygo

/*
 * TinyGo stub: TinyGo's `os` package is missing several APIs used by the
 * generic syscall namespace (Chmod, etc). Browser-targeted wasm has no
 * meaningful filesystem anyway, so this just installs a no-op namespace.
 */

package rt

func installSyscallNS() {}

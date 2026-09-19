//go:build !(linux || darwin || freebsd || openbsd || netbsd || dragonfly || illumos)

// bench-baton relies on flock(2) and process groups for its leases. The
// supported list is explicit rather than `unix`: aix and solaris are `unix`
// but lack syscall.Flock, so they get this stub instead of a failed build.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "bench-baton: not supported on this platform (needs flock and process groups)")
	os.Exit(2)
}

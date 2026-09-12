//go:build !unix

// bench-baton relies on flock(2) for its leases; on platforms without it
// the command explains itself instead of failing to build ./...
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "bench-baton: only supported on Unix-like systems (needs flock)")
	os.Exit(2)
}

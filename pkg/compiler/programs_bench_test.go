package compiler

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// BenchmarkPrograms runs the fixed Clojure programs under benchmark/*.clj as
// end-to-end workloads. Each program is compiled ONCE (outside the timer); the
// timed loop only executes the compiled chunk, so this measures execution speed
// of the compiled output — the axis the IR-optimization pipeline affects. The
// .clj files are committed and identical across versions, so this is a stable
// cross-version comparison (unlike BenchmarkClojureTestSuite, whose workload is
// the evolving test suite).
func BenchmarkPrograms(b *testing.B) {
	dir := filepath.Join("..", "..", "benchmark")
	files, err := filepath.Glob(filepath.Join(dir, "*.clj"))
	if err != nil || len(files) == 0 {
		b.Skipf("no benchmark programs found in %s", dir)
	}
	sort.Strings(files)
	for _, f := range files {
		name := strings.TrimSuffix(filepath.Base(f), ".clj")
		src, err := os.ReadFile(f)
		if err != nil {
			b.Fatalf("read %s: %v", f, err)
		}
		b.Run(name, func(b *testing.B) {
			consts := vm.NewConsts()
			ctx := NewCompiler(consts, rt.NS(rt.NameCoreNS))
			chunk, _, err := ctx.CompileMultiple(strings.NewReader(string(src)))
			if err != nil {
				b.Skipf("compile %s: %v", name, err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fr := vm.NewFrame(chunk, nil)
				fr.Run()
			}
		})
	}
}

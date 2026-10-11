package bytecode

import (
	"bytes"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// deferTestBundle builds a two-namespace bundle by hand: each namespace's main
// chunk loads its own function constants; the shared keyword constant sits in
// core's range and is used by both. The returned ranges mirror what lgbgen
// records (the pool span each compilation appended).
//
// Layout (chunk indices follow ModuleBuilder: main chunks first, then the
// function chunks in const order):
//
//	consts: 0 :shared, 1 fn a1, 2 fn a2, 3 fn b1, 4 :b-only
//	chunks: 0 main-core, 1 main-b, 2 a1, 3 a2, 4 b1
//
// padB1 appends that many never-executed NOOPs after b1's RETURN, to grow b's
// code without changing anything else.
func deferTestBundle(t *testing.T, shareFuncAcrossNS bool, padB1 int) (data []byte, ranges map[string]ConstRange) {
	t.Helper()
	pool := vm.NewConsts()
	mkFn := func(name string, retConst int, maxStack int) *vm.Func {
		c := vm.NewCodeChunk(pool)
		c.Append(vm.OP_LOAD_CONST, int32(retConst), vm.OP_RETURN)
		if name == "b1" {
			for i := 0; i < padB1; i++ {
				c.Append(vm.OP_NOOP)
			}
		}
		c.SetMaxStack(maxStack) // distinct per fn: the Module re-encode matches chunks structurally
		c.AddSourceInfoAt(0, vm.SourceInfo{File: name + ".lg", Line: 1, Column: 1, EndLine: 1, EndColumn: 5})
		c.AddLocalVar(0, "x-"+name)
		f := vm.MakeFunc(0, false, c)
		f.SetName(name)
		return f
	}
	shared := pool.Intern(vm.Keyword("shared")) // 0
	a1 := pool.Intern(mkFn("a1", shared, 1))    // 1
	a2 := pool.Intern(mkFn("a2", shared, 2))    // 2
	bOnlyIdx := 4
	b1 := pool.Intern(mkFn("b1", bOnlyIdx, 3)) // 3
	pool.Intern(vm.Keyword("b-only"))          // 4

	mainA := vm.NewCodeChunk(pool)
	mainA.Append(vm.OP_LOAD_CONST, int32(a1), vm.OP_POP, vm.OP_LOAD_CONST, int32(a2), vm.OP_POP, vm.OP_LOAD_CONST, int32(shared), vm.OP_RETURN)
	mainA.SetMaxStack(1)
	mainB := vm.NewCodeChunk(pool)
	if shareFuncAcrossNS {
		mainB.Append(vm.OP_LOAD_CONST, int32(a1), vm.OP_POP)
	}
	mainB.Append(vm.OP_LOAD_CONST, int32(b1), vm.OP_POP, vm.OP_LOAD_CONST, int32(shared), vm.OP_RETURN)
	mainB.SetMaxStack(1)
	mainB.AddLocalVar(0, "x-main-b")

	var buf bytes.Buffer
	ranges = map[string]ConstRange{"core": {Lo: 0, Hi: 3}, "b": {Lo: 3, Hi: 5}}
	if err := EncodeBundleDeferrable(&buf, pool, map[string]*vm.CodeChunk{"core": mainA, "b": mainB}, []string{"core", "b"}, ranges, false); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes(), ranges
}

func TestNSRangesProvenFromReachability(t *testing.T) {
	data, _ := deferTestBundle(t, false, 0)
	m, err := Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if m.Flags&FlagNSRanges == 0 {
		t.Fatal("FlagNSRanges not set")
	}
	// :shared (const 0) is reached by both namespaces, so it is in no run;
	// :b-only (const 4) is reached by b alone and is deferred with b.
	want := []NSRange{
		{Name: "b", MainChunk: 1, ChunkLo: 4, ChunkHi: 5, ConstRuns: []ConstRange{{Lo: 3, Hi: 5}}},
		{Name: "core", MainChunk: 0, ChunkLo: 2, ChunkHi: 4, ConstRuns: []ConstRange{{Lo: 1, Hi: 3}}},
	}
	if !reflect.DeepEqual(m.NSRanges, want) {
		t.Fatalf("ranges = %+v, want %+v", m.NSRanges, want)
	}

	// A re-encode of the decoded Module keeps the section byte-for-byte.
	var again bytes.Buffer
	if err := Encode(&again, m); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(again.Bytes(), data) {
		t.Fatal("re-encoding a module with namespace ranges changed its bytes")
	}
}

func TestNSRangesOmitNamespaceWhoseFunctionIsShared(t *testing.T) {
	data, _ := deferTestBundle(t, true, 0)
	m, err := Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	// b loads a1, so core's function a1 is reached by two namespaces: core
	// must not be deferrable. b's own function b1 is still reached only by b.
	if len(m.NSRanges) != 1 || m.NSRanges[0].Name != "b" {
		t.Fatalf("ranges = %+v, want only b", m.NSRanges)
	}
	// And deferring core is silently a no-op: core decodes eagerly.
	unit, err := DecodeBundle(data, nil, func(string) bool { return true })
	if err != nil {
		t.Fatal(err)
	}
	if unit.NSChunks["core"] == nil || unit.DeferredNS("core") {
		t.Fatal("core should have decoded eagerly")
	}
	if unit.NSChunks["b"] != nil || !unit.DeferredNS("b") {
		t.Fatal("b should be deferred")
	}
}

// TestDecodeBundleDefersAndMaterializesIdentically decodes the same bundle
// eagerly and with b deferred, then materializes b and compares every chunk,
// function constant, source map and local-var table.
func TestDecodeBundleDefersAndMaterializesIdentically(t *testing.T) {
	data, _ := deferTestBundle(t, false, 0)
	eager, err := DecodeToExecUnitBytes(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	SetDecodeStatsEnabled(true)
	defer SetDecodeStatsEnabled(false)
	ResetDecodeStats()
	lazy, err := DecodeBundle(data, nil, func(name string) bool { return name == "b" })
	if err != nil {
		t.Fatal(err)
	}
	if got := SnapshotDecodeStats().Chunks; got != 3 {
		t.Fatalf("deferred decode materialized %d chunks, want 3 (main-core, a1, a2)", got)
	}
	if lazy.NSChunks["b"] != nil || !lazy.DeferredNS("b") {
		t.Fatal("b is not deferred")
	}
	if len(lazy.NSOrder) != 2 || lazy.NSOrder[1] != "b" {
		t.Fatalf("NSOrder = %v", lazy.NSOrder)
	}
	pool := lazy.Consts.AllValues()
	if len(pool) != 5 {
		t.Fatalf("pool has %d slots, want 5 (deferred slots are reserved)", len(pool))
	}
	if pool[3] != nil {
		t.Fatalf("b1's slot holds %v before materialization", pool[3])
	}
	if pool[4] != nil {
		t.Fatalf("b's own constant :b-only was decoded eagerly: %v", pool[4])
	}
	if pool[0] != vm.Keyword("shared") {
		t.Fatalf("the shared constant was not decoded eagerly: %v", pool[0])
	}
	if _, ok := pool[1].(*vm.Func); !ok {
		t.Fatalf("a1 was not decoded: %v", pool[1])
	}

	ResetDecodeStats()
	ch, err := lazy.NSChunk("b")
	if err != nil {
		t.Fatal(err)
	}
	if got := SnapshotDecodeStats().Chunks; got != 2 {
		t.Fatalf("materializing b decoded %d chunks, want 2 (main-b, b1)", got)
	}
	if lazy.DeferredNS("b") {
		t.Fatal("b still reported deferred")
	}
	compareChunk(t, "main-b", eager.NSChunks["b"], ch)
	ef, lf := eager.Consts.AllValues()[3].(*vm.Func), pool[3].(*vm.Func)
	if ef.FuncName() != lf.FuncName() || ef.Arity() != lf.Arity() || ef.IsVariadic() != lf.IsVariadic() {
		t.Fatalf("b1 differs: eager %s/%d lazy %s/%d", ef.FuncName(), ef.Arity(), lf.FuncName(), lf.Arity())
	}
	compareChunk(t, "b1", ef.Chunk(), lf.Chunk())
	if lazy.Consts.IsDeferred(3) || pool[4] != vm.Keyword("b-only") {
		t.Fatalf("b's constants not filled after materialization: %v %v", pool[3], pool[4])
	}

	// Materialization is idempotent and returns the same chunk.
	again, err := lazy.NSChunk("b")
	if err != nil || again != ch {
		t.Fatalf("second NSChunk: %v, %v", again, err)
	}
	// And the eager namespace is served from NSChunks as before.
	if got, err := lazy.NSChunk("core"); err != nil || got != lazy.NSChunks["core"] || got != lazy.MainChunk {
		t.Fatalf("NSChunk(core) = %v, %v", got, err)
	}
	if got, err := lazy.NSChunk("nope"); err != nil || got != nil {
		t.Fatalf("NSChunk(unknown) = %v, %v", got, err)
	}

	// Running b's main chunk yields the shared keyword, as the eager one does.
	for _, u := range []*ExecUnit{eager, lazy} {
		mb, _ := u.NSChunk("b")
		f := vm.NewFrame(mb, nil)
		out, err := f.RunProtected()
		vm.ReleaseFrame(f)
		if err != nil || out != vm.Keyword("shared") {
			t.Fatalf("running main-b: %v, %v", out, err)
		}
	}
}

func compareChunk(t *testing.T, what string, e, l *vm.CodeChunk) {
	t.Helper()
	if e == nil || l == nil {
		t.Fatalf("%s: nil chunk (eager=%v lazy=%v)", what, e != nil, l != nil)
	}
	ec, lc := e.Code(), l.Code()
	if len(ec) != len(lc) {
		t.Fatalf("%s: code length eager=%d lazy=%d", what, len(ec), len(lc))
	}
	for i := range ec {
		if ec[i] != lc[i] {
			t.Fatalf("%s: code[%d] eager=%d lazy=%d", what, i, ec[i], lc[i])
		}
	}
	if e.MaxStack() != l.MaxStack() {
		t.Fatalf("%s: max stack eager=%d lazy=%d", what, e.MaxStack(), l.MaxStack())
	}
	for ip := 0; ip <= len(ec); ip++ {
		es, ls := e.LookupSource(ip), l.LookupSource(ip)
		if (es == nil) != (ls == nil) || (es != nil && *es != *ls) {
			t.Fatalf("%s: source at %d eager=%v lazy=%v", what, ip, es, ls)
		}
	}
	elv, llv := e.LocalVars(), l.LocalVars()
	if len(elv) != len(llv) {
		t.Fatalf("%s: local vars eager=%v lazy=%v", what, elv, llv)
	}
	for i := range elv {
		if elv[i] != llv[i] {
			t.Fatalf("%s: local var[%d] eager=%v lazy=%v", what, i, elv[i], llv[i])
		}
	}
}

func TestDecodeBundleMaterializesOnceUnderConcurrentRequires(t *testing.T) {
	data, _ := deferTestBundle(t, false, 0)
	unit, err := DecodeBundle(data, nil, func(name string) bool { return true })
	if err != nil {
		t.Fatal(err)
	}
	if unit.DeferredNS("core") || unit.MainChunk == nil {
		t.Fatal("the unit's main chunk (core, chunk 0) must never be deferred")
	}
	const n = 16
	got := make([]*vm.CodeChunk, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch, err := unit.NSChunk("b")
			if err != nil {
				t.Error(err)
			}
			got[i] = ch
		}(i)
	}
	wg.Wait()
	for i := 1; i < n; i++ {
		if got[i] != got[0] {
			t.Fatalf("concurrent NSChunk returned different chunks")
		}
	}
	f := unit.Consts.AllValues()[3].(*vm.Func)
	if f.Chunk() == nil {
		t.Fatal("b1 has no chunk")
	}
}

// TestDecodeBundleAllocationsIgnoreDeferredCode is the property the startup
// gate rests on: growing a deferred namespace's code does not move what the
// deferring decode allocates. Both bundles decode identically apart from the
// padded (never executed) instructions in b's function.
func TestDecodeBundleAllocationsIgnoreDeferredCode(t *testing.T) {
	data, _ := deferTestBundle(t, false, 0)
	grown, _ := deferTestBundle(t, false, 4096)
	if len(grown) <= len(data)+4096*4 {
		t.Fatalf("padding did not grow the bundle: %d -> %d bytes", len(data), len(grown))
	}
	deferB := func(name string) bool { return name == "b" }
	decode := func(d []byte) func() {
		return func() {
			if _, err := DecodeBundle(d, nil, deferB); err != nil {
				t.Fatal(err)
			}
		}
	}
	small := testing.AllocsPerRun(50, decode(data))
	big := testing.AllocsPerRun(50, decode(grown))
	if small != big {
		t.Fatalf("deferring decode allocs moved with deferred code: %v -> %v", small, big)
	}
	eagerSmall := testing.AllocsPerRun(50, func() { DecodeToExecUnitBytes(data, nil) })
	eagerBig := testing.AllocsPerRun(50, func() { DecodeToExecUnitBytes(grown, nil) })
	t.Logf("deferring: %v allocs both; eager: %v -> %v", small, eagerSmall, eagerBig)
}

// TestSkipValueV2MatchesRead pins skipValueV2 to readValueV2: over every
// constant of the synthetic bundle and of the real core bundle, skipping a
// record must consume exactly the bytes reading it does.
func TestSkipValueV2MatchesRead(t *testing.T) {
	synthetic, _ := deferTestBundle(t, false, 0)
	bundles := map[string][]byte{"synthetic": synthetic}
	if core, err := os.ReadFile("../rt/core_compiled.lgb"); err == nil {
		bundles["core"] = core
	}
	for name, data := range bundles {
		m, err := Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		// Re-encode uncompressed, then walk to the consts section by reading
		// the sections before it the way the decoder does.
		var buf bytes.Buffer
		m.Flags &^= FlagCompressed
		m.Version = uncompressedFormatVersion
		if err := Encode(&buf, m); err != nil {
			t.Fatal(err)
		}
		mk := func() *decoder {
			d := &decoder{r: NewReaderBytes(buf.Bytes())}
			if _, _, err := d.readHeader(); err != nil {
				t.Fatal(err)
			}
			d.flags = m.Flags
			if err := d.readCapabilities(); err != nil {
				t.Fatal(err)
			}
			if d.strings, err = d.readStringTable(); err != nil {
				t.Fatal(err)
			}
			if d.flags&FlagNSRanges != 0 {
				if _, err := d.readNSRanges(); err != nil {
					t.Fatal(err)
				}
			}
			if err := d.readLiveChunks(vm.NewConsts(), nil); err != nil {
				t.Fatal(err)
			}
			if _, err := d.r.ReadVarint(); err != nil { // const count
				t.Fatal(err)
			}
			if d.flags&FlagConstsBase != 0 {
				if _, err := d.r.ReadVarint(); err != nil {
					t.Fatal(err)
				}
			}
			return d
		}
		rd, sk := mk(), mk()
		tags := map[byte]int{}
		for i := range m.Consts {
			tags[tagOf(t, rd)]++
			if _, err := rd.readValueV2(); err != nil {
				t.Fatalf("%s const[%d]: read: %v", name, i, err)
			}
			if err := sk.skipValueV2(); err != nil {
				t.Fatalf("%s const[%d]: skip: %v", name, i, err)
			}
			if rd.r.Offset() != sk.r.Offset() {
				t.Fatalf("%s const[%d] (%v): read ended at %d, skip at %d", name, i, m.Consts[i], rd.r.Offset(), sk.r.Offset())
			}
		}
		t.Logf("%s: %d consts, tags %v", name, len(m.Consts), tags)
	}
}

// tagOf peeks the tag byte of the next record without consuming it.
func tagOf(t *testing.T, d *decoder) byte {
	t.Helper()
	b, err := d.r.r.Peek(1)
	if err != nil {
		t.Fatal(err)
	}
	return b[0] & tagIDMask
}

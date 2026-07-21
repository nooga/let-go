package bytecode

import (
	"bytes"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// --- Encoding primitive roundtrips ---

func TestVarintRoundtrip(t *testing.T) {
	cases := []uint64{0, 1, 127, 128, 255, 256, 16383, 16384, math.MaxUint64}
	for _, tc := range cases {
		var buf bytes.Buffer
		w := NewWriter(&buf)
		if err := w.WriteVarint(tc); err != nil {
			t.Fatalf("WriteVarint(%d): %v", tc, err)
		}
		w.Flush()
		r := NewReader(&buf)
		got, err := r.ReadVarint()
		if err != nil {
			t.Fatalf("ReadVarint for %d: %v", tc, err)
		}
		if got != tc {
			t.Errorf("varint roundtrip: got %d, want %d", got, tc)
		}
	}
}

func TestSvarintRoundtrip(t *testing.T) {
	cases := []int64{0, 1, -1, 127, -128, math.MaxInt64, math.MinInt64}
	for _, tc := range cases {
		var buf bytes.Buffer
		w := NewWriter(&buf)
		if err := w.WriteSvarint(tc); err != nil {
			t.Fatalf("WriteSvarint(%d): %v", tc, err)
		}
		w.Flush()
		r := NewReader(&buf)
		got, err := r.ReadSvarint()
		if err != nil {
			t.Fatalf("ReadSvarint for %d: %v", tc, err)
		}
		if got != tc {
			t.Errorf("svarint roundtrip: got %d, want %d", got, tc)
		}
	}
}

func TestFloat64Roundtrip(t *testing.T) {
	cases := []float64{0.0, 1.0, -1.0, 3.14159, math.MaxFloat64, math.SmallestNonzeroFloat64, math.Inf(1), math.Inf(-1)}
	for _, tc := range cases {
		var buf bytes.Buffer
		w := NewWriter(&buf)
		if err := w.WriteFloat64(tc); err != nil {
			t.Fatalf("WriteFloat64(%v): %v", tc, err)
		}
		w.Flush()
		r := NewReader(&buf)
		got, err := r.ReadFloat64()
		if err != nil {
			t.Fatalf("ReadFloat64 for %v: %v", tc, err)
		}
		if got != tc {
			t.Errorf("float64 roundtrip: got %v, want %v", got, tc)
		}
	}
	// NaN special case
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.WriteFloat64(math.NaN())
	w.Flush()
	r := NewReader(&buf)
	got, _ := r.ReadFloat64()
	if !math.IsNaN(got) {
		t.Errorf("NaN roundtrip: got %v", got)
	}
}

func TestInt32Roundtrip(t *testing.T) {
	cases := []int32{0, 1, -1, math.MaxInt32, math.MinInt32}
	for _, tc := range cases {
		var buf bytes.Buffer
		w := NewWriter(&buf)
		if err := w.WriteInt32(tc); err != nil {
			t.Fatalf("WriteInt32(%d): %v", tc, err)
		}
		w.Flush()
		r := NewReader(&buf)
		got, err := r.ReadInt32()
		if err != nil {
			t.Fatalf("ReadInt32 for %d: %v", tc, err)
		}
		if got != tc {
			t.Errorf("int32 roundtrip: got %d, want %d", got, tc)
		}
	}
}

func TestUint16Roundtrip(t *testing.T) {
	cases := []uint16{0, 1, 65535}
	for _, tc := range cases {
		var buf bytes.Buffer
		w := NewWriter(&buf)
		if err := w.WriteUint16(tc); err != nil {
			t.Fatalf("WriteUint16(%d): %v", tc, err)
		}
		w.Flush()
		r := NewReader(&buf)
		got, err := r.ReadUint16()
		if err != nil {
			t.Fatalf("ReadUint16 for %d: %v", tc, err)
		}
		if got != tc {
			t.Errorf("uint16 roundtrip: got %d, want %d", got, tc)
		}
	}
}

// --- Value roundtrips ---

func roundtripModule(t *testing.T, consts []vm.Value, chunks []*ChunkData) *Module {
	t.Helper()
	sb := NewModuleBuilder()
	for _, v := range consts {
		sb.internStringsForValue(v)
	}
	if chunks != nil {
		for _, ch := range chunks {
			for _, sm := range ch.SourceMap {
				sb.internString(sm.File)
			}
		}
	}
	m := &Module{
		Version: FormatVersion,
		Flags:   0,
		Consts:  consts,
		Chunks:  chunks,
		Strings: sb.strings,
	}

	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	return decoded
}

func roundtripValue(t *testing.T, v vm.Value) vm.Value {
	t.Helper()
	m := roundtripModule(t, []vm.Value{v}, nil)
	if len(m.Consts) != 1 {
		t.Fatalf("expected 1 const, got %d", len(m.Consts))
	}
	return m.Consts[0]
}

func TestNilRoundtrip(t *testing.T) {
	got := roundtripValue(t, vm.NIL)
	if got != vm.NIL {
		t.Errorf("expected NIL, got %v", got)
	}
}

func TestBoolRoundtrip(t *testing.T) {
	got := roundtripValue(t, vm.TRUE)
	if got != vm.TRUE {
		t.Errorf("expected TRUE, got %v", got)
	}
	got = roundtripValue(t, vm.FALSE)
	if got != vm.FALSE {
		t.Errorf("expected FALSE, got %v", got)
	}
}

func TestIntRoundtrip(t *testing.T) {
	cases := []int{0, 42, -1, math.MaxInt64, math.MinInt64}
	for _, tc := range cases {
		got := roundtripValue(t, vm.Int(tc))
		if got != vm.Int(tc) {
			t.Errorf("Int roundtrip: got %v, want %v", got, vm.Int(tc))
		}
	}
}

func TestFloatRoundtrip(t *testing.T) {
	cases := []float64{0.0, 3.14, -1.5}
	for _, tc := range cases {
		got := roundtripValue(t, vm.Float(tc))
		if got != vm.Float(tc) {
			t.Errorf("Float roundtrip: got %v, want %v", got, vm.Float(tc))
		}
	}
}

func TestStringRoundtrip(t *testing.T) {
	cases := []string{"", "hello", "日本語", "with\x00null"}
	for _, tc := range cases {
		got := roundtripValue(t, vm.String(tc))
		if got != vm.String(tc) {
			t.Errorf("String roundtrip: got %v, want %v", got, vm.String(tc))
		}
	}
}

func TestKeywordRoundtrip(t *testing.T) {
	cases := []string{"foo", "ns/bar"}
	for _, tc := range cases {
		got := roundtripValue(t, vm.Keyword(tc))
		if got != vm.Keyword(tc) {
			t.Errorf("Keyword roundtrip: got %v, want %v", got, vm.Keyword(tc))
		}
	}
}

func TestSymbolRoundtrip(t *testing.T) {
	cases := []string{"baz", "ns/quux"}
	for _, tc := range cases {
		got := roundtripValue(t, vm.Symbol(tc))
		if got != vm.Symbol(tc) {
			t.Errorf("Symbol roundtrip: got %v, want %v", got, vm.Symbol(tc))
		}
	}
}

func TestCharRoundtrip(t *testing.T) {
	cases := []rune{'a', 'λ', 0}
	for _, tc := range cases {
		got := roundtripValue(t, vm.Char(tc))
		if got != vm.Char(tc) {
			t.Errorf("Char roundtrip: got %v, want %v", got, vm.Char(tc))
		}
	}
}

func TestBigIntRoundtrip(t *testing.T) {
	cases := []*big.Int{
		big.NewInt(0),
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
		new(big.Int).Neg(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)),
	}
	for _, tc := range cases {
		got := roundtripValue(t, vm.NewBigInt(tc))
		bi, ok := got.(*vm.BigInt)
		if !ok {
			t.Fatalf("expected *BigInt, got %T", got)
		}
		if bi.Val().Cmp(tc) != 0 {
			t.Errorf("BigInt roundtrip: got %v, want %v", bi.Val(), tc)
		}
	}
}

func TestVoidRoundtrip(t *testing.T) {
	got := roundtripValue(t, vm.VOID)
	if got != vm.VOID {
		t.Errorf("expected VOID, got %v", got)
	}
}

func TestUUIDRoundtrip(t *testing.T) {
	cases := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"00000000-0000-0000-0000-000000000000",
		"ffffffff-ffff-ffff-ffff-ffffffffffff",
	}
	for _, tc := range cases {
		u := vm.NewUUID(tc)
		got := roundtripValue(t, u)
		gu, ok := got.(*vm.UUID)
		if !ok {
			t.Fatalf("expected *UUID, got %T", got)
		}
		if s := gu.Unbox().(string); s != tc {
			t.Errorf("UUID roundtrip: got %q, want %q", s, tc)
		}
	}
}

// Uppercase input must canonicalise to lowercase through the reader (ParseUUID),
// and that lowercase form must survive the bytecode round-trip unchanged.
func TestUUIDRoundtripCanonicalisesUppercase(t *testing.T) {
	u := vm.ParseUUID("550E8400-E29B-41D4-A716-446655440000")
	if u == nil {
		t.Fatal("ParseUUID returned nil for uppercase input")
	}
	got := roundtripValue(t, u)
	gu, ok := got.(*vm.UUID)
	if !ok {
		t.Fatalf("expected *UUID, got %T", got)
	}
	want := "550e8400-e29b-41d4-a716-446655440000"
	if s := gu.Unbox().(string); s != want {
		t.Errorf("UUID canonicalisation: got %q, want %q", s, want)
	}
}

func TestEmptyListRoundtrip(t *testing.T) {
	got := roundtripValue(t, vm.EmptyList)
	if got != vm.EmptyList {
		t.Errorf("expected EmptyList, got %v", got)
	}
}

// --- Func roundtrip ---

func TestFuncRoundtrip(t *testing.T) {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)

	fn := vm.MakeFunc(2, false, chunk)
	fn.SetName("test-fn")

	// Build module using the builder
	b := NewModuleBuilder()
	b.AddChunk(chunk)
	b.AddConst(fn)
	m := b.Build()

	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if len(decoded.Consts) != 1 {
		t.Fatalf("expected 1 const, got %d", len(decoded.Consts))
	}
	gotFn, ok := decoded.Consts[0].(*vm.Func)
	if !ok {
		t.Fatalf("expected *Func, got %T", decoded.Consts[0])
	}
	if gotFn.FuncName() != "test-fn" {
		t.Errorf("name: got %q, want %q", gotFn.FuncName(), "test-fn")
	}
	if gotFn.Arity() != 2 {
		t.Errorf("arity: got %d, want 2", gotFn.Arity())
	}
	if gotFn.IsVariadic() {
		t.Error("expected non-variadic")
	}
	// Check bytecode
	gotCode := gotFn.Chunk().Code()
	wantCode := []int32{vm.OP_LOAD_CONST, 0, vm.OP_RETURN}
	if len(gotCode) != len(wantCode) {
		t.Fatalf("code length: got %d, want %d", len(gotCode), len(wantCode))
	}
	for i := range wantCode {
		if gotCode[i] != wantCode[i] {
			t.Errorf("code[%d]: got %d, want %d", i, gotCode[i], wantCode[i])
		}
	}
}

func TestFuncVariadicRoundtrip(t *testing.T) {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_RETURN)
	chunk.SetMaxStack(1)

	fn := vm.MakeFunc(3, true, chunk)
	fn.SetName("varfn")

	b := NewModuleBuilder()
	b.AddChunk(chunk)
	b.AddConst(fn)
	m := b.Build()

	var buf bytes.Buffer
	Encode(&buf, m)
	decoded, _ := Decode(&buf)
	gotFn := decoded.Consts[0].(*vm.Func)
	if !gotFn.IsVariadic() {
		t.Error("expected variadic")
	}
	if gotFn.Arity() != 3 {
		t.Errorf("arity: got %d, want 3", gotFn.Arity())
	}
}

// --- Var ref roundtrip ---

func TestVarRefRoundtrip(t *testing.T) {
	v := vm.NewVar(nil, "user", "my-var")
	v.SetRoot(vm.Int(42))

	b := NewModuleBuilder()
	b.AddConst(v)
	m := b.Build()

	var buf bytes.Buffer
	Encode(&buf, m)

	resolved := vm.NewVar(nil, "user", "my-var")
	resolved.SetRoot(vm.Int(99))

	decoded, err := DecodeWithResolver(&buf, func(ns, name string) *vm.Var {
		if ns == "user" && name == "my-var" {
			return resolved
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	gotVar, ok := decoded.Consts[0].(*vm.Var)
	if !ok {
		t.Fatalf("expected *Var, got %T", decoded.Consts[0])
	}
	if gotVar != resolved {
		t.Error("expected resolved var to be the same pointer")
	}
}

func TestVarRefWithoutResolver(t *testing.T) {
	v := vm.NewVar(nil, "core", "println")
	b := NewModuleBuilder()
	b.AddConst(v)
	m := b.Build()

	var buf bytes.Buffer
	Encode(&buf, m)
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	gotVar := decoded.Consts[0].(*vm.Var)
	if gotVar.NS() != "core" || gotVar.VarName() != "println" {
		t.Errorf("var: got %s/%s, want core/println", gotVar.NS(), gotVar.VarName())
	}
}

// --- Collection roundtrips ---

func TestListRoundtrip(t *testing.T) {
	items := []vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)}
	l, _ := vm.ListType.Box(items)

	got := roundtripValue(t, l.(*vm.List))
	gotList, ok := got.(*vm.List)
	if !ok {
		t.Fatalf("expected *List, got %T", got)
	}
	if gotList.RawCount() != 3 {
		t.Errorf("count: got %d, want 3", gotList.RawCount())
	}
	if gotList.First() != vm.Int(1) {
		t.Errorf("first: got %v, want 1", gotList.First())
	}
}

func TestMixedList(t *testing.T) {
	items := []vm.Value{vm.Keyword("a"), vm.String("b"), vm.Int(42)}
	l, _ := vm.ListType.Box(items)
	got := roundtripValue(t, l.(*vm.List))
	gotList := got.(*vm.List)
	if gotList.First() != vm.Keyword("a") {
		t.Errorf("first: got %v", gotList.First())
	}
}

func TestVectorRoundtrip(t *testing.T) {
	v := vm.ArrayVector{vm.Int(1), vm.Int(2), vm.Int(3)}
	got := roundtripValue(t, v)
	gotVec, ok := got.(vm.ArrayVector)
	if !ok {
		t.Fatalf("expected ArrayVector, got %T", got)
	}
	if len(gotVec) != 3 {
		t.Errorf("length: got %d, want 3", len(gotVec))
	}
	for i, want := range []vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)} {
		if gotVec[i] != want {
			t.Errorf("vec[%d]: got %v, want %v", i, gotVec[i], want)
		}
	}
}

func TestEmptyVectorRoundtrip(t *testing.T) {
	v := vm.ArrayVector{}
	got := roundtripValue(t, v)
	gotVec := got.(vm.ArrayVector)
	if len(gotVec) != 0 {
		t.Errorf("expected empty vector, got length %d", len(gotVec))
	}
}

func TestMapRoundtrip(t *testing.T) {
	m := vm.EmptyPersistentMap
	m = m.Assoc(vm.Keyword("x"), vm.Int(1)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("y"), vm.Int(2)).(*vm.PersistentMap)

	got := roundtripValue(t, m)
	gotMap, ok := got.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected *PersistentMap, got %T", got)
	}
	if gotMap.RawCount() != 2 {
		t.Errorf("count: got %d, want 2", gotMap.RawCount())
	}
	if gotMap.ValueAt(vm.Keyword("x")) != vm.Int(1) {
		t.Errorf(":x got %v, want 1", gotMap.ValueAt(vm.Keyword("x")))
	}
	if gotMap.ValueAt(vm.Keyword("y")) != vm.Int(2) {
		t.Errorf(":y got %v, want 2", gotMap.ValueAt(vm.Keyword("y")))
	}
}

func TestEmptyMapRoundtrip(t *testing.T) {
	got := roundtripValue(t, vm.EmptyPersistentMap)
	gotMap := got.(*vm.PersistentMap)
	if gotMap.RawCount() != 0 {
		t.Errorf("expected empty map, got count %d", gotMap.RawCount())
	}
}

func TestNestedMapRoundtrip(t *testing.T) {
	inner := vm.EmptyPersistentMap
	inner = inner.Assoc(vm.Keyword("b"), vm.Int(1)).(*vm.PersistentMap)
	outer := vm.EmptyPersistentMap
	outer = outer.Assoc(vm.Keyword("a"), inner).(*vm.PersistentMap)

	got := roundtripValue(t, outer)
	gotMap := got.(*vm.PersistentMap)
	innerGot := gotMap.ValueAt(vm.Keyword("a")).(*vm.PersistentMap)
	if innerGot.ValueAt(vm.Keyword("b")) != vm.Int(1) {
		t.Error("nested map value mismatch")
	}
}

func TestLargeMapRoundtrip(t *testing.T) {
	m := vm.EmptyPersistentMap
	for i := range 100 {
		key := string(rune('a'+i/26)) + string(rune('a'+i%26))
		m = m.Assoc(vm.Keyword(key), vm.Int(i)).(*vm.PersistentMap)
	}

	got := roundtripValue(t, m)
	gotMap, ok := got.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected *PersistentMap, got %T", got)
	}
	if gotMap.RawCount() != 100 {
		t.Errorf("count: got %d, want 100", gotMap.RawCount())
	}
	for i := range 100 {
		key := string(rune('a'+i/26)) + string(rune('a'+i%26))
		kw := vm.Keyword(key)
		if gotMap.ValueAt(kw) != vm.Int(i) {
			t.Errorf("key %v: got %v, want %d", kw, gotMap.ValueAt(kw), i)
		}
	}
}

func TestSetRoundtrip(t *testing.T) {
	s := vm.NewPersistentSet([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)})
	got := roundtripValue(t, s)
	gotSet, ok := got.(*vm.PersistentSet)
	if !ok {
		t.Fatalf("expected *PersistentSet, got %T", got)
	}
	if gotSet.RawCount() != 3 {
		t.Errorf("count: got %d, want 3", gotSet.RawCount())
	}
	if gotSet.Contains(vm.Int(1)) != vm.TRUE {
		t.Error("set should contain 1")
	}
}

func TestLargeVectorRoundtrip(t *testing.T) {
	vec := make(vm.ArrayVector, 100)
	for i := range 100 {
		vec[i] = vm.Int(i)
	}
	got := roundtripValue(t, vec)
	gotVec, ok := got.(vm.ArrayVector)
	if !ok {
		t.Fatalf("expected ArrayVector, got %T", got)
	}
	if len(gotVec) != 100 {
		t.Errorf("length: got %d, want 100", len(gotVec))
	}
	for i := range 100 {
		if gotVec[i] != vm.Int(i) {
			t.Errorf("vec[%d]: got %v, want %d", i, gotVec[i], i)
		}
	}
}

func TestLargeSetRoundtrip(t *testing.T) {
	items := make([]vm.Value, 100)
	for i := range 100 {
		items[i] = vm.Int(i)
	}
	s := vm.NewPersistentSet(items)
	got := roundtripValue(t, s)
	gotSet, ok := got.(*vm.PersistentSet)
	if !ok {
		t.Fatalf("expected *PersistentSet, got %T", got)
	}
	if gotSet.RawCount() != 100 {
		t.Errorf("count: got %d, want 100", gotSet.RawCount())
	}
	for i := range 100 {
		if gotSet.Contains(vm.Int(i)) != vm.TRUE {
			t.Errorf("set should contain %d", i)
		}
	}
}

func TestEmptySetRoundtrip(t *testing.T) {
	s := vm.NewPersistentSet(nil)
	got := roundtripValue(t, s)
	gotSet := got.(*vm.PersistentSet)
	if gotSet.RawCount() != 0 {
		t.Errorf("expected empty set, got count %d", gotSet.RawCount())
	}
}

func TestNestedCollections(t *testing.T) {
	// Vector containing a map containing a list
	inner, _ := vm.ListType.Box([]vm.Value{vm.Int(1), vm.Int(2)})
	m := vm.EmptyPersistentMap
	m = m.Assoc(vm.Keyword("items"), inner).(*vm.PersistentMap)
	v := vm.ArrayVector{m, vm.String("end")}

	got := roundtripValue(t, v)
	gotVec := got.(vm.ArrayVector)
	gotMap := gotVec[0].(*vm.PersistentMap)
	gotList := gotMap.ValueAt(vm.Keyword("items")).(*vm.List)
	if gotList.First() != vm.Int(1) {
		t.Error("nested value mismatch")
	}
	if gotVec[1] != vm.String("end") {
		t.Errorf("got %v, want end", gotVec[1])
	}
}

// --- RecordType roundtrip ---

func TestRecordTypeRoundtrip(t *testing.T) {
	rt := vm.NewRecordType("MyRecord", []vm.Keyword{"x", "y"})
	got := roundtripValue(t, rt)
	gotRT, ok := got.(*vm.RecordType)
	if !ok {
		t.Fatalf("expected *RecordType, got %T", got)
	}
	if gotRT.TypeName() != "MyRecord" {
		t.Errorf("name: got %q, want %q", gotRT.TypeName(), "MyRecord")
	}
	fields := gotRT.Fields()
	if len(fields) != 2 || fields[0] != "x" || fields[1] != "y" {
		t.Errorf("fields: got %v", fields)
	}
}

// --- Record roundtrip ---

func TestRecordRoundtrip(t *testing.T) {
	rt := vm.NewRecordType("Point", []vm.Keyword{"x", "y"})
	data := vm.EmptyPersistentMap
	data = data.Assoc(vm.Keyword("x"), vm.Int(10)).(*vm.PersistentMap)
	data = data.Assoc(vm.Keyword("y"), vm.Int(20)).(*vm.PersistentMap)
	rec := vm.NewRecord(rt, data)

	got := roundtripValue(t, rec)
	gotRec, ok := got.(*vm.Record)
	if !ok {
		t.Fatalf("expected *Record, got %T", got)
	}
	if gotRec.ValueAt(vm.Keyword("x")) != vm.Int(10) {
		t.Errorf(":x got %v, want 10", gotRec.ValueAt(vm.Keyword("x")))
	}
	if gotRec.ValueAt(vm.Keyword("y")) != vm.Int(20) {
		t.Errorf(":y got %v, want 20", gotRec.ValueAt(vm.Keyword("y")))
	}
}

func TestRecordWithExtraFields(t *testing.T) {
	rt := vm.NewRecordType("Ext", []vm.Keyword{"a"})
	data := vm.EmptyPersistentMap
	data = data.Assoc(vm.Keyword("a"), vm.Int(1)).(*vm.PersistentMap)
	data = data.Assoc(vm.Keyword("extra"), vm.String("bonus")).(*vm.PersistentMap)
	rec := vm.NewRecord(rt, data)

	got := roundtripValue(t, rec).(*vm.Record)
	if got.ValueAt(vm.Keyword("a")) != vm.Int(1) {
		t.Error(":a mismatch")
	}
	if got.ValueAt(vm.Keyword("extra")) != vm.String("bonus") {
		t.Error(":extra mismatch")
	}
}

func TestEmptyRecordRoundtrip(t *testing.T) {
	rt := vm.NewRecordType("Empty", []vm.Keyword{"a"})
	rec := vm.NewRecord(rt, vm.EmptyPersistentMap)

	got := roundtripValue(t, rec).(*vm.Record)
	if got.ValueAt(vm.Keyword("a")) != vm.NIL {
		t.Errorf(":a should be nil, got %v", got.ValueAt(vm.Keyword("a")))
	}
}

// --- Regex roundtrip ---

func TestRegexRoundtrip(t *testing.T) {
	v, err := vm.NewRegex(`foo\d+`)
	if err != nil {
		t.Fatal(err)
	}
	got := roundtripValue(t, v)
	gotRe, ok := got.(*vm.Regex)
	if !ok {
		t.Fatalf("expected *Regex, got %T", got)
	}
	if gotRe.Pattern() != `foo\d+` {
		t.Errorf("pattern: got %q, want %q", gotRe.Pattern(), `foo\d+`)
	}
}

// A terminal-lookahead pattern only compiles through vm.NewRegex's
// compatibility fallback — raw regexp.Compile rejects it. The decoder must
// reconstruct through the same path or the bundle round-trip fails with
// "invalid or unsupported Perl syntax" (review finding on #494).
func TestRegexLookaheadRoundtrip(t *testing.T) {
	const pat = `(\w)-(?=\w)`
	v, err := vm.NewRegex(pat)
	if err != nil {
		t.Fatal(err)
	}
	got := roundtripValue(t, v)
	gotRe, ok := got.(*vm.Regex)
	if !ok {
		t.Fatalf("expected *Regex, got %T", got)
	}
	if gotRe.Pattern() != pat {
		t.Errorf("pattern: got %q, want %q", gotRe.Pattern(), pat)
	}
	// The decoded regex must MATCH with lookahead semantics (visible match
	// excludes the assertion), not just carry the pattern string.
	if m := gotRe.FindStringSubmatch("a-b"); len(m) == 0 || m[0] != "a-" {
		t.Errorf("decoded lookahead match: got %v, want visible match %q", m, "a-")
	}
}

// --- Atom roundtrip ---

func TestAtomRoundtrip(t *testing.T) {
	a := vm.NewAtom(vm.Int(42))
	got := roundtripValue(t, a)
	gotAtom, ok := got.(*vm.Atom)
	if !ok {
		t.Fatalf("expected *Atom, got %T", got)
	}
	if gotAtom.Deref() != vm.Int(42) {
		t.Errorf("deref: got %v, want 42", gotAtom.Deref())
	}
}

func TestAtomWithMapValue(t *testing.T) {
	m := vm.EmptyPersistentMap
	m = m.Assoc(vm.Keyword("state"), vm.String("ready")).(*vm.PersistentMap)
	a := vm.NewAtom(m)

	got := roundtripValue(t, a).(*vm.Atom)
	gotMap := got.Deref().(*vm.PersistentMap)
	if gotMap.ValueAt(vm.Keyword("state")) != vm.String("ready") {
		t.Error("atom map value mismatch")
	}
}

// --- Full module roundtrip ---

func TestFullModuleRoundtrip(t *testing.T) {
	consts := vm.NewConsts()
	chunk1 := vm.NewCodeChunk(consts)
	chunk1.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk1.SetMaxStack(2)

	chunk2 := vm.NewCodeChunk(consts)
	chunk2.Append(vm.OP_LOAD_ARG, 0, vm.OP_RETURN)
	chunk2.SetMaxStack(1)

	chunk3 := vm.NewCodeChunk(consts)
	chunk3.Append(vm.OP_ADD, vm.OP_RETURN)
	chunk3.SetMaxStack(3)

	fn1 := vm.MakeFunc(0, false, chunk1)
	fn1.SetName("foo")
	fn2 := vm.MakeFunc(1, false, chunk2)
	fn2.SetName("bar")
	fn3 := vm.MakeFunc(2, true, chunk3)
	fn3.SetName("foo") // same name as fn1 — tests string dedup

	b := NewModuleBuilder()
	b.AddChunk(chunk1)
	b.AddChunk(chunk2)
	b.AddChunk(chunk3)

	allConsts := []vm.Value{
		vm.NIL, vm.TRUE, vm.FALSE,
		vm.Int(0), vm.Int(42), vm.Int(-1),
		vm.Float(3.14), vm.Float(0.0),
		vm.String("hello"), vm.String("world"),
		vm.Keyword("foo"), vm.Keyword("bar"),
		vm.Symbol("baz"), vm.Symbol("quux"),
		vm.Char('x'), vm.Char('λ'),
		vm.VOID,
		vm.EmptyList,
		fn1, fn2, fn3,
	}
	for _, v := range allConsts {
		b.AddConst(v)
	}
	m := b.Build()

	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(decoded.Consts) != len(allConsts) {
		t.Fatalf("const count: got %d, want %d", len(decoded.Consts), len(allConsts))
	}
	if len(decoded.Chunks) != 3 {
		t.Fatalf("chunk count: got %d, want 3", len(decoded.Chunks))
	}
}

func TestV2Roundtrip(t *testing.T) {
	// Explicit v2 marker: verify the encoder writes version 2 and the decoder
	// correctly reads it back on a module containing common scalar types.
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(2)

	b := NewModuleBuilder()
	b.AddChunk(chunk)

	allConsts := []vm.Value{
		vm.NIL, vm.TRUE, vm.FALSE,
		vm.Int(0), vm.Int(42), vm.Int(-1),
		vm.Float(3.14), vm.Float(0.0),
		vm.String("hello"), vm.String("world"),
		vm.Keyword("foo"), vm.Keyword("bar"),
		vm.Symbol("baz"), vm.Symbol("quux"),
		vm.Char('x'), vm.Char('λ'),
		vm.VOID,
		vm.EmptyList,
	}
	for _, v := range allConsts {
		b.AddConst(v)
	}
	m := b.Build()
	if m.Version != 2 {
		t.Fatalf("expected version 2, got %d", m.Version)
	}

	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if decoded.Version != 2 {
		t.Fatalf("decoded version: got %d, want 2", decoded.Version)
	}
	if len(decoded.Consts) != len(allConsts) {
		t.Fatalf("const count: got %d, want %d", len(decoded.Consts), len(allConsts))
	}
	for i, want := range allConsts {
		got := decoded.Consts[i]
		if got != want {
			t.Errorf("const[%d]: got %v, want %v", i, got, want)
		}
	}
}

func TestStringDedup(t *testing.T) {
	b := NewModuleBuilder()
	b.AddConst(vm.Keyword("foo"))
	b.AddConst(vm.Symbol("foo"))
	fn := vm.MakeFunc(0, false, vm.NewCodeChunk(vm.NewConsts()))
	fn.SetName("foo")
	b.AddChunk(fn.Chunk())
	b.AddConst(fn)
	m := b.Build()

	count := 0
	for _, s := range m.Strings {
		if s == "foo" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 'foo' once in string table, got %d", count)
	}
}

// --- Source map roundtrip ---

func TestSourceMapRoundtrip(t *testing.T) {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_LOAD_CONST, 1, vm.OP_ADD, vm.OP_RETURN)
	chunk.SetMaxStack(2)
	chunk.AddSourceInfo(vm.SourceInfo{File: "test.lg", Line: 0, Column: 0, EndLine: 0, EndColumn: 5})
	chunk.AddSourceInfo(vm.SourceInfo{File: "test.lg", Line: 1, Column: 2, EndLine: 1, EndColumn: 10})
	chunk.AddSourceInfo(vm.SourceInfo{File: "other.lg", Line: 5, Column: 0, EndLine: 5, EndColumn: 20})

	b := NewModuleBuilder()
	b.AddChunk(chunk)
	m := b.Build()

	var buf bytes.Buffer
	Encode(&buf, m)
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(decoded.Chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(decoded.Chunks))
	}
	sm := decoded.Chunks[0].SourceMap
	if len(sm) != 3 {
		t.Fatalf("expected 3 source map entries, got %d", len(sm))
	}
	if sm[0].File != "test.lg" || sm[0].Line != 0 || sm[0].Column != 0 {
		t.Errorf("sm[0]: got %+v", sm[0])
	}
	if sm[2].File != "other.lg" || sm[2].Line != 5 || sm[2].EndColumn != 20 {
		t.Errorf("sm[2]: got %+v", sm[2])
	}
}

// --- Edge cases ---

func TestEmptyModule(t *testing.T) {
	m := &Module{Version: FormatVersion, Flags: 0}
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(decoded.Strings) != 0 {
		t.Errorf("expected 0 strings, got %d", len(decoded.Strings))
	}
	if len(decoded.Chunks) != 0 {
		t.Errorf("expected 0 chunks, got %d", len(decoded.Chunks))
	}
	if len(decoded.Consts) != 0 {
		t.Errorf("expected 0 consts, got %d", len(decoded.Consts))
	}
}

func TestLargeString(t *testing.T) {
	s := make([]byte, 65536)
	for i := range s {
		s[i] = byte('a' + (i % 26))
	}
	got := roundtripValue(t, vm.String(s))
	if got != vm.String(s) {
		t.Error("large string roundtrip failed")
	}
}

func TestManyConsts(t *testing.T) {
	consts := make([]vm.Value, 1000)
	for i := range consts {
		consts[i] = vm.Int(i)
	}
	m := roundtripModule(t, consts, nil)
	if len(m.Consts) != 1000 {
		t.Errorf("expected 1000 consts, got %d", len(m.Consts))
	}
	for i, v := range m.Consts {
		if v != vm.Int(i) {
			t.Errorf("const[%d]: got %v, want %d", i, v, i)
			break
		}
	}
}

// --- Error handling ---

func TestTruncatedHeader(t *testing.T) {
	_, err := Decode(bytes.NewReader([]byte("LG")))
	if err == nil {
		t.Error("expected error for truncated header")
	}
}

func TestInvalidMagic(t *testing.T) {
	_, err := Decode(bytes.NewReader([]byte("XXXX\x01\x00\x00\x00")))
	if err == nil {
		t.Error("expected error for invalid magic")
	}
}

func TestUnknownTag(t *testing.T) {
	// Build a valid module but with a bogus tag in the const pool
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.WriteBytes(Magic[:])
	w.WriteUint16(FormatVersion)
	w.WriteUint16(0)
	w.WriteVarint(0)  // 0 strings
	w.WriteVarint(0)  // 0 chunks
	w.WriteVarint(1)  // 1 const
	w.WriteByte(0xFF) // unknown tag
	w.Flush()

	_, err := Decode(&buf)
	if err == nil {
		t.Error("expected error for unknown tag")
	}
}

func TestVarintOverflow(t *testing.T) {
	// 11 continuation bytes (exceeds 10 byte limit)
	data := make([]byte, 11)
	for i := range data {
		data[i] = 0x80
	}
	r := NewReader(bytes.NewReader(data))
	_, err := r.ReadVarint()
	if err == nil {
		t.Error("expected error for varint overflow")
	}
}

func TestChunkIndexOutOfRange(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.WriteBytes(Magic[:])
	w.WriteUint16(FormatVersion)
	w.WriteUint16(0)
	w.WriteVarint(1) // 1 string
	w.WriteVarint(0) // string ""
	w.WriteVarint(0) // 0 chunks
	w.WriteVarint(1) // 1 const
	w.WriteByte(TagFunc)
	w.WriteVarint(99) // chunk index out of range
	w.WriteVarint(0)  // arity
	w.WriteByte(0)    // not variadic
	w.WriteVarint(0)  // name ref
	w.Flush()

	_, err := Decode(&buf)
	if err == nil {
		t.Error("expected error for chunk index out of range")
	}
}

func TestUnknownTagError(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.WriteBytes(Magic[:])
	w.WriteUint16(2)  // version 2
	w.WriteUint16(0)  // flags
	w.WriteVarint(0)  // 0 strings
	w.WriteVarint(0)  // 0 chunks
	w.WriteVarint(1)  // 1 const
	w.WriteByte(0xFF) // unknown tag ID 0x3F with version 3 (0b11_111111)
	w.Flush()

	_, err := Decode(&buf)
	if err == nil {
		t.Fatal("expected error for unknown tag")
	}
	msg := err.Error()
	if !strings.Contains(msg, "unknown tag") {
		t.Fatalf("error message should mention 'unknown tag', got: %s", msg)
	}
}

func TestV1DecodeStillWorks(t *testing.T) {
	// Hand-craft a minimal v1 module and ensure the frozen v1 path decodes it.
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.WriteBytes(Magic[:])
	w.WriteUint16(1) // version 1
	w.WriteUint16(0) // flags
	w.WriteVarint(1) // 1 string
	w.WriteVarint(5) // length "hello"
	w.WriteBytes([]byte("hello"))
	w.WriteVarint(0)       // 0 chunks
	w.WriteVarint(1)       // 1 const
	w.WriteByte(TagString) // TagString in v1 is 0x05, same byte as v2
	w.WriteVarint(0)       // string ref 0
	w.Flush()

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("v1 decode failed: %v", err)
	}
	if decoded.Version != 1 {
		t.Fatalf("version: got %d, want 1", decoded.Version)
	}
	if len(decoded.Consts) != 1 {
		t.Fatalf("expected 1 const, got %d", len(decoded.Consts))
	}
	if decoded.Consts[0] != vm.String("hello") {
		t.Errorf("const: got %v, want hello", decoded.Consts[0])
	}
}

func TestUnsupportedCapabilityMask(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.WriteBytes(Magic[:])
	w.WriteUint16(2)
	w.WriteUint16(FlagCapabilities)
	w.WriteUint32(0x00000002) // bit 1 set — no such capability defined
	w.Flush()

	_, err := Decode(&buf)
	if err == nil {
		t.Fatal("expected error for unsupported capability mask")
	}
	if !strings.Contains(err.Error(), "unsupported capability mask") {
		t.Fatalf("expected 'unsupported capability mask' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "recompile it with a matching lg") {
		t.Fatalf("expected the recompile hint in error, got: %v", err)
	}
}

func TestSupportedCapabilityMask(t *testing.T) {
	// capability mask 0x00000000 should be accepted even when FlagCapabilities is set
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.WriteBytes(Magic[:])
	w.WriteUint16(2)
	w.WriteUint16(FlagCapabilities)
	w.WriteUint32(0) // no capabilities required
	w.WriteVarint(0) // 0 strings
	w.WriteVarint(0) // 0 chunks
	w.WriteVarint(0) // 0 consts
	w.Flush()

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.Version != 2 {
		t.Fatalf("version: got %d, want 2", decoded.Version)
	}
}

// --- Determinism ---

// makeBundleChunks builds N tiny NS chunks for determinism testing.
// Names are picked so Go's randomized map iteration would surface
// quickly across runs.
func makeBundleChunks(consts *vm.Consts) map[string]*vm.CodeChunk {
	names := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}
	out := make(map[string]*vm.CodeChunk, len(names))
	for i, n := range names {
		ch := vm.NewCodeChunk(consts)
		ch.Append(vm.OP_LOAD_CONST, int32(i), vm.OP_RETURN)
		ch.SetMaxStack(1)
		out[n] = ch
	}
	return out
}

// TestEncodeBundleOrderedIsDeterministic ensures EncodeBundleOrdered
// produces byte-identical output across invocations with the same order.
func TestEncodeBundleOrderedIsDeterministic(t *testing.T) {
	order := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}
	var first []byte
	for i := range 8 {
		consts := vm.NewConsts()
		chunks := makeBundleChunks(consts)
		var buf bytes.Buffer
		if err := EncodeBundleOrdered(&buf, consts, chunks, order); err != nil {
			t.Fatalf("encode %d: %v", i, err)
		}
		if i == 0 {
			first = buf.Bytes()
			continue
		}
		if !bytes.Equal(first, buf.Bytes()) {
			t.Fatalf("EncodeBundleOrdered run %d differs from run 0", i)
		}
	}
}

// TestEncodeBundleIsDeterministic ensures the map-input EncodeBundle also
// produces byte-identical output across runs (it sorts NS names internally).
func TestEncodeBundleIsDeterministic(t *testing.T) {
	var first []byte
	for i := range 8 {
		consts := vm.NewConsts()
		chunks := makeBundleChunks(consts)
		var buf bytes.Buffer
		if err := EncodeBundle(&buf, consts, chunks); err != nil {
			t.Fatalf("encode %d: %v", i, err)
		}
		if i == 0 {
			first = buf.Bytes()
			continue
		}
		if !bytes.Equal(first, buf.Bytes()) {
			t.Fatalf("EncodeBundle run %d differs from run 0", i)
		}
	}
}

// --- Benchmarks ---

func BenchmarkVarintRoundtrip(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	r := NewReader(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.WriteVarint(uint64(i))
		w.Flush()
		r.ReadVarint()
	}
}

func BenchmarkEncodeModule(b *testing.B) {
	mb := NewModuleBuilder()
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)
	mb.AddChunk(chunk)

	for i := range 100 {
		mb.AddConst(vm.Int(i))
	}
	for range 20 {
		mb.AddConst(vm.Keyword("key"))
	}
	m := mb.Build()

	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		Encode(&buf, m)
	}
}

func BenchmarkDecodeModule(b *testing.B) {
	mb := NewModuleBuilder()
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)
	mb.AddChunk(chunk)

	for i := range 100 {
		mb.AddConst(vm.Int(i))
	}
	for range 20 {
		mb.AddConst(vm.Keyword("key"))
	}
	m := mb.Build()

	var buf bytes.Buffer
	Encode(&buf, m)
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(bytes.NewReader(data))
	}
}

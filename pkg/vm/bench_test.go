package vm

import (
	"testing"
)

// ============================================================================
// Frame & Stack operations — measures raw VM dispatch overhead
// ============================================================================

// BenchmarkFrameDispatch measures the cost of the bytecode dispatch loop
// with a tight loop of LOAD_CONST + POP.
func BenchmarkFrameDispatch(b *testing.B) {
	consts := NewConsts()
	idx := consts.Intern(Int(42))

	// Build: LOAD_CONST idx; POP; LOAD_CONST idx; POP; ... LOAD_CONST idx; RETURN
	// 100 iterations of load+pop, then one final load+return
	chunk := NewCodeChunk(consts)
	iterations := 100
	for i := 0; i < iterations; i++ {
		chunk.Append(OP_LOAD_CONST | (int32(i%8) << 16))
		chunk.Append32(idx)
		chunk.Append(OP_POP)
	}
	chunk.Append(OP_LOAD_CONST)
	chunk.Append32(idx)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := NewFrame(chunk, nil)
		f.Run()
	}
}

// BenchmarkFrameAlloc measures frame allocation cost
func BenchmarkFrameAlloc(b *testing.B) {
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_CONST)
	chunk.Append32(consts.Intern(NIL))
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)

	args := []Value{Int(1), Int(2), Int(3)}

	b.Run("NoArgs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewFrame(chunk, nil)
		}
	})
	b.Run("3Args", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewFrame(chunk, args)
		}
	})
}

// ============================================================================
// Function call overhead — measures invoke path costs
// ============================================================================

// BenchmarkFuncInvoke benchmarks calling compiled functions
func BenchmarkFuncInvoke(b *testing.B) {
	consts := NewConsts()

	// identity function: LOAD_ARG 0; RETURN
	fnChunk := NewCodeChunk(consts)
	fnChunk.Append(OP_LOAD_ARG)
	fnChunk.Append32(0)
	fnChunk.Append(OP_RETURN)
	fnChunk.SetMaxStack(4)
	fn := MakeFunc(1, false, fnChunk)

	args := []Value{Int(42)}

	b.Run("Direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fn.Invoke(args)
		}
	})

	// Same via closure
	closure := fn.MakeClosure()
	b.Run("Closure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			closure.Invoke(args)
		}
	})

	// NativeFn
	nativeFn, _ := NativeFnType.Wrap(func(vs []Value) (Value, error) {
		return vs[0], nil
	})
	b.Run("Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nativeFn.(Fn).Invoke(args)
		}
	})
}

// BenchmarkVariadicInvoke benchmarks variadic function call overhead
func BenchmarkVariadicInvoke(b *testing.B) {
	consts := NewConsts()

	// variadic identity: LOAD_ARG 0; RETURN
	fnChunk := NewCodeChunk(consts)
	fnChunk.Append(OP_LOAD_ARG)
	fnChunk.Append32(0)
	fnChunk.Append(OP_RETURN)
	fnChunk.SetMaxStack(4)
	fn := MakeFunc(1, true, fnChunk)

	b.Run("1Arg", func(b *testing.B) {
		args := []Value{Int(1)}
		for i := 0; i < b.N; i++ {
			fn.Invoke(args)
		}
	})
	b.Run("5Args", func(b *testing.B) {
		args := []Value{Int(1), Int(2), Int(3), Int(4), Int(5)}
		for i := 0; i < b.N; i++ {
			fn.Invoke(args)
		}
	})
}

// BenchmarkMultiArity benchmarks multi-arity dispatch
func BenchmarkMultiArity(b *testing.B) {
	consts := NewConsts()

	// Two arities: 1-arg and 2-arg, both just return first arg
	fn1Chunk := NewCodeChunk(consts)
	fn1Chunk.Append(OP_LOAD_ARG)
	fn1Chunk.Append32(0)
	fn1Chunk.Append(OP_RETURN)
	fn1Chunk.SetMaxStack(4)
	fn1 := MakeFunc(1, false, fn1Chunk)

	fn2Chunk := NewCodeChunk(consts)
	fn2Chunk.Append(OP_LOAD_ARG)
	fn2Chunk.Append32(0)
	fn2Chunk.Append(OP_RETURN)
	fn2Chunk.SetMaxStack(4)
	fn2 := MakeFunc(2, false, fn2Chunk)

	ma, _ := makeMultiArity([]Value{fn1, fn2})

	b.Run("1Arg", func(b *testing.B) {
		args := []Value{Int(42)}
		for i := 0; i < b.N; i++ {
			ma.Invoke(args)
		}
	})
	b.Run("2Args", func(b *testing.B) {
		args := []Value{Int(42), Int(99)}
		for i := 0; i < b.N; i++ {
			ma.Invoke(args)
		}
	})
}

// ============================================================================
// Stack push/pop micro-benchmarks
// ============================================================================

func BenchmarkStackOps(b *testing.B) {
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_NOOP)
	chunk.SetMaxStack(64)

	val := Int(42)

	b.Run("Push", func(b *testing.B) {
		f := NewFrame(chunk, nil)
		for i := 0; i < b.N; i++ {
			f.sp = 0
			f.push(val)
		}
	})
	b.Run("PushPop", func(b *testing.B) {
		f := NewFrame(chunk, nil)
		for i := 0; i < b.N; i++ {
			f.sp = 0
			f.push(val)
			f.pop()
		}
	})
}

// ============================================================================
// Seq iteration — measures cost of walking different seq types
// ============================================================================

func BenchmarkSeqIteration(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		values := make([]Value, size)
		for i := range values {
			values[i] = Int(i)
		}

		list, _ := ListType.Box(values)
		vec := NewArrayVector(values)

		b.Run("List/"+itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for s := Seq(list.(*List)); s != nil; s = s.Next() {
					_ = s.First()
				}
			}
		})

		b.Run("ArrayVector/"+itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sq := vec.(Sequable).Seq()
				for s := sq; s != nil; s = s.Next() {
					_ = s.First()
				}
			}
		})

		if size <= 1000 {
			pv := NewPersistentVector(values)
			b.Run("PersistentVector/"+itoa(size), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					sq := pv.(Sequable).Seq()
					for s := sq; s != nil; s = s.Next() {
						_ = s.First()
					}
				}
			})
		}
	}
}

// ============================================================================
// Collection operations
// ============================================================================

func BenchmarkConj(b *testing.B) {
	b.Run("List", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var c Collection = EmptyList
			for j := 0; j < 100; j++ {
				c = c.Conj(Int(j))
			}
		}
	})
	b.Run("ArrayVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var c Collection = ArrayVector{}
			for j := 0; j < 100; j++ {
				c = c.Conj(Int(j))
			}
		}
	})
	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var c Collection = make(Map)
			for j := 0; j < 100; j++ {
				c = c.Conj(ArrayVector{Int(j), Int(j * 10)})
			}
		}
	})
}

func BenchmarkMapAssoc(b *testing.B) {
	sizes := []int{10, 100, 1000}
	for _, size := range sizes {
		// Build initial map
		m := make(Map, size)
		for i := 0; i < size; i++ {
			m[Int(i)] = Int(i * 10)
		}
		b.Run("Assoc/"+itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m.Assoc(Int(0), Int(999))
			}
		})
		b.Run("Dissoc/"+itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m.Dissoc(Int(0))
			}
		})
		b.Run("Lookup/"+itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m.ValueAt(Int(size / 2))
			}
		})
	}
}

// ============================================================================
// IsTruthy — hot path in every branch
// ============================================================================

func BenchmarkIsTruthy(b *testing.B) {
	vals := []struct {
		name string
		v    Value
	}{
		{"true", TRUE},
		{"false", FALSE},
		{"nil", NIL},
		{"int", Int(42)},
		{"string", String("hello")},
	}
	for _, tc := range vals {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsTruthy(tc.v)
			}
		})
	}
}

// ============================================================================
// Cons cell creation — hot in lazy seq chains
// ============================================================================

func BenchmarkConsCreation(b *testing.B) {
	val := Int(42)
	b.Run("NewCons", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewCons(val, nil)
		}
	})
	b.Run("ListCons", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			EmptyList.Cons(val)
		}
	})
}

// ============================================================================
// LazySeq realization
// ============================================================================

func BenchmarkLazySeqRealize(b *testing.B) {
	// Build a chain of N lazy seqs
	makeLazy := func(n int) *LazySeq {
		var build func(i int) *LazySeq
		build = func(i int) *LazySeq {
			if i >= n {
				return NewLazySeq(nil)
			}
			captured := i
			thunk, _ := NativeFnType.Wrap(func(_ []Value) (Value, error) {
				return NewCons(Int(captured), build(captured+1)), nil
			})
			return NewLazySeq(thunk.(Fn))
		}
		return build(0)
	}

	sizes := []int{10, 100, 1000}
	for _, size := range sizes {
		b.Run("Realize/"+itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ls := makeLazy(size)
				for s := Seq(ls); s != nil; s = s.Next() {
					_ = s.First()
				}
			}
		})
	}
}

func itoa(i int) string {
	switch {
	case i < 10:
		return string(rune('0' + i))
	case i < 100:
		return string([]rune{rune('0' + i/10), rune('0' + i%10)})
	case i < 1000:
		return string([]rune{rune('0' + i/100), rune('0' + (i/10)%10), rune('0' + i%10)})
	default:
		return string([]rune{rune('0' + i/10000), rune('0' + (i/1000)%10), rune('0' + (i/100)%10), rune('0' + (i/10)%10), rune('0' + i%10)})
	}
}

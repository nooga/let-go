---
status: active
last-verified: 2026-11-19
authoritative-for:
  - type-lattice-kind-mask
  - numeric-track
human-verified:
---

# The type lattice's kind mask and the numeric track

`ir.lattice` (`pkg/rt/core/ir/lattice.lg`) keeps two kinds of type fact apart. This note records that split, decided in June 2026 when the bitset lattice landed, so the inference work that depends on it has one place to cite.

## The kind mask is a pure set

The kind mask has one bit per kind (`:nil :true :false :bool :int :float :string :char`, plus `:any` and `:unknown`); `:number` is `:int` ∪ `:float`. Join is bitwise OR, and the subset order is the lattice order: a ⊑ b exactly when a's bits are a subset of b's. The mask's integer value is not a total order, and nothing reads it as one.

The mask path is proven equal to the keyword path: `scripts/ti-bitset-proptest.lg` checks mask join against `type-join` over all 26,244 non-constant pairs, and the normalize fast path was checked against the slow path over all 511 kind subsets when it landed (the Go test `lattice.lg` cites for that, `TestZZNormalizeMaskEquiv`, is not on main).

Kinds without a bit (`:ratio`, collection kinds, `[:dtype T]`, arrays) travel as keyword facts on the slow path, unchanged.

## Numbers beyond the kind live in the numeric track

Numeric values (const-zero, ranges) stay out of the kind mask and live in a separate numeric track. Tower bounds live there too: a bound (the lowest set bit) means something only inside the numeric category, so it never enters the kind-lattice join. Those two points are the June decision.

The encoding #1139 proposes for the tower: one bit per representation in widening order (fixnum widths i8, i16, i32, i64; BigInt; Ratio; BigDecimal; float widths f32, f64). A set bit means a value may have that representation. The lowest set bit is the lower bound, the narrowest representation that may occur. The representation that holds every value is the highest set bit, but only within one chain where each wider bit holds every narrower one (the fixnum widths, or the float widths); a set that spans chains, such as fixnum and float, or anything with BigInt, Ratio or BigDecimal beside another bit, has no single unboxed representation.

The numeric track is not built. The branch that started it (`numeric-handlers-integration`, type-specialized numeric handlers via gogen) was retired without landing.

## What depends on it

The execution-form spec's bounds for a parameter, `join(caller argument types) ⊑ T ⊑ meet(usage constraints)` (`docs/superpowers/specs/2026-10-03-execution-form-design.md`), need the usage meet to say which numbers an op accepts. Arithmetic and comparison accept the whole tower; bit ops and checked indexes accept fixnums up to i64; a float op accepts the float bits. With the numeric track those are per-op tower sets, and the meet over a param's uses picks its representation: fixnum bits only unbox to the highest set width (the narrowest that holds every value), float bits only to the highest set float width, and a set spanning chains stays boxed or takes the spec's guarded typed entry. #1139 tracks that work.

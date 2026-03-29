package vm

import "unsafe"

// Consts is a constant pool with O(1) indexed access and O(1) amortized interning.
// Uses a custom open-addressing hash map (not Go's built-in map) so that all Value
// types — including slices like ArrayVector — can be interned without panicking.
//
// Supports optional layering: a child pool references a parent. Indices 0..base-1
// are in the parent, base.. are in the child. Only child entries are returned by
// Values() for serialization.
type Consts struct {
	parent *Consts // nil for root pool
	base   int     // parent.count() at creation; child indices start here

	consts []Value // this layer's values

	// Open-addressing hash map for dedup (linear probing, power-of-2 sizing)
	buckets []constBucket
	size    int    // occupied buckets
	mask    uint32 // len(buckets) - 1
}

type constBucket struct {
	hash  uint32
	index int   // global index
	value Value // nil = empty slot
}

// NewConsts creates a root constant pool.
func NewConsts() *Consts {
	return &Consts{
		consts:  make([]Value, 0, 64),
		buckets: make([]constBucket, 16),
		mask:    15,
	}
}

// NewChildConsts creates a child pool layered on top of parent.
// Parent values are accessible by index but not serialized with this pool.
func NewChildConsts(parent *Consts) *Consts {
	return &Consts{
		parent:  parent,
		base:    parent.count(),
		consts:  make([]Value, 0, 32),
		buckets: make([]constBucket, 16),
		mask:    15,
	}
}

// Intern returns the index for v, adding it if not already present.
// Checks parent pool first (if any) to reuse existing indices.
func (c *Consts) Intern(v Value) int {
	h := constHash(v)

	// Check parent first
	if c.parent != nil {
		if idx, ok := c.parent.lookup(h, v); ok {
			return idx
		}
	}

	// Check this layer
	if idx, ok := c.lookup(h, v); ok {
		return idx
	}

	// Add new entry
	idx := c.base + len(c.consts)
	c.consts = append(c.consts, v)
	c.insertBucket(h, idx, v)
	return idx
}

// Append adds a value at the next index without deduplication.
// Used when loading pre-compiled bytecode where indices must match the original.
func (c *Consts) Append(v Value) int {
	idx := c.base + len(c.consts)
	c.consts = append(c.consts, v)
	h := constHash(v)
	c.insertBucket(h, idx, v)
	return idx
}

func (c *Consts) get(i int) Value {
	if i >= c.base {
		return c.consts[i-c.base]
	}
	return c.parent.consts[i]
}

func (c *Consts) count() int {
	return c.base + len(c.consts)
}

// Values returns only this layer's values (for serialization).
func (c *Consts) Values() []Value {
	return c.consts
}

// AllValues returns all values including parent.
func (c *Consts) AllValues() []Value {
	if c.parent == nil {
		return c.consts
	}
	all := make([]Value, 0, c.count())
	all = append(all, c.parent.consts...)
	all = append(all, c.consts...)
	return all
}

// Base returns the starting global index for this layer.
func (c *Consts) Base() int {
	return c.base
}

// Parent returns the parent pool, or nil.
func (c *Consts) Parent() *Consts {
	return c.parent
}

// --- open-addressing hash map ---

func (c *Consts) lookup(h uint32, v Value) (int, bool) {
	if c.mask == 0 {
		return 0, false
	}
	pos := h & c.mask
	for {
		e := &c.buckets[pos]
		if e.value == nil {
			return 0, false
		}
		if e.hash == h && constEqual(e.value, v) {
			return e.index, true
		}
		pos = (pos + 1) & c.mask
	}
}

func (c *Consts) insertBucket(h uint32, idx int, v Value) {
	// Grow at 75% load
	if c.size*4 >= int(c.mask+1)*3 {
		c.grow()
	}
	pos := h & c.mask
	for {
		e := &c.buckets[pos]
		if e.value == nil {
			e.hash = h
			e.index = idx
			e.value = v
			c.size++
			return
		}
		pos = (pos + 1) & c.mask
	}
}

func (c *Consts) grow() {
	newSize := (c.mask + 1) * 2
	if newSize < 16 {
		newSize = 16
	}
	old := c.buckets
	c.buckets = make([]constBucket, newSize)
	c.mask = uint32(newSize - 1)
	c.size = 0
	for _, e := range old {
		if e.value != nil {
			c.insertBucket(e.hash, e.index, e.value)
		}
	}
}

// --- hash and equality for const pool ---

// constHash computes a hash for deduplication in the const pool.
// Pointer-identity types hash by address; structural types use hashValue.
func constHash(v Value) uint32 {
	switch x := v.(type) {
	case *Func:
		return hashUint64(uint64(uintptr(unsafe.Pointer(x))))
	case *Var:
		return hashUint64(uint64(uintptr(unsafe.Pointer(x))))
	case *BigInt:
		return hashUint64(uint64(uintptr(unsafe.Pointer(x))))
	case *Record:
		return hashUint64(uint64(uintptr(unsafe.Pointer(x))))
	case *RecordType:
		return hashUint64(uint64(uintptr(unsafe.Pointer(x))))
	case *Regex:
		return hashUint64(uint64(uintptr(unsafe.Pointer(x))))
	case *Atom:
		return hashUint64(uint64(uintptr(unsafe.Pointer(x))))
	default:
		return hashValue(v)
	}
}

// constEqual tests equality for const pool deduplication.
// Pointer-identity types use pointer equality; value types use same-type
// structural comparison (never merging across types like Int/Float/BigInt).
func constEqual(a, b Value) bool {
	switch a.(type) {
	case *Func, *Var, *BigInt, *Record, *RecordType, *Regex, *Atom:
		return a == b
	default:
		// Only consider equal if same concrete type
		if a.Type() != b.Type() {
			return false
		}
		return valueEquiv(a, b)
	}
}

package vm

import (
	"fmt"
	"reflect"
	"time"
)

type theInstantType struct{}

func (t *theInstantType) String() string  { return t.Name() }
func (t *theInstantType) Type() ValueType { return TypeType }
func (t *theInstantType) Unbox() any      { return reflect.TypeFor[*theInstantType]() }
func (t *theInstantType) Name() string    { return "let-go.lang.Instant" }
func (t *theInstantType) Box(bare any) (Value, error) {
	switch v := bare.(type) {
	case string:
		i := ParseInstant(v)
		if i == nil {
			return NIL, fmt.Errorf("invalid instant: %s", v)
		}
		return i, nil
	case time.Time:
		// RFC3339(Nano) requires a 4-digit year (0000..9999), so a
		// time.Time outside that range cannot round-trip through
		// Val()/ParseInstant — which would silently break AOT bytecode
		// serialisation. Refuse at the host-interop boundary rather than
		// producing a Val() string our own decoder can't read back.
		// Year 0 is allowed because RFC3339 permits "0000-..." and
		// time.Parse(RFC3339Nano) accepts it on the way back in.
		if y := v.Year(); y < 0 || y > 9999 {
			return NIL, fmt.Errorf("time.Time year %d outside RFC3339 range (0..9999)", y)
		}
		return NewInstant(v), nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", InstantType)
}

// InstantType is the type of Instant values.
var InstantType *theInstantType = &theInstantType{}

// Instant wraps a time.Time, normalised to UTC at millisecond precision.
//
// Precision is millisecond to match Clojure JVM (java.util.Date) and
// ClojureScript (js/Date), so values round-trip cleanly through any
// Clojure tool reading the same RFC3339 string.
type Instant struct {
	t time.Time
}

// ParseInstant parses an RFC3339 (or RFC3339Nano) timestamp. Returns nil if invalid.
// Sub-millisecond precision in the input is truncated.
func ParseInstant(s string) *Instant {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return nil
	}
	return &Instant{t: t.UTC().Truncate(time.Millisecond)}
}

// NewInstant constructs an Instant from a time.Time, normalising to UTC
// at millisecond precision.
//
// The caller is responsible for ensuring the year is in 0..9999 — values
// outside that range will not round-trip through Val()/ParseInstant and
// will silently break AOT bytecode serialisation. The only value-system
// entry point that takes a host time.Time is InstantType.Box, which
// enforces this bound; direct Go callers must do the same.
func NewInstant(t time.Time) *Instant {
	return &Instant{t: t.UTC().Truncate(time.Millisecond)}
}

func (i *Instant) Type() ValueType { return InstantType }
func (i *Instant) Unbox() any      { return i.t }
func (i *Instant) String() string  { return "#inst \"" + i.t.Format(time.RFC3339Nano) + "\"" }

// Time returns the underlying time.Time (UTC, millisecond precision).
func (i *Instant) Time() time.Time { return i.t }

// Val returns the canonical RFC3339Nano string form (no #inst wrapping).
func (i *Instant) Val() string { return i.t.Format(time.RFC3339Nano) }

// Hash implements Hashable. Grounded in epoch milliseconds (the actual
// value) rather than the formatted string so hash identity is independent
// of presentation choices (RFC3339Nano trims trailing zeros, etc.).
func (i *Instant) Hash() uint32 { return hashUint64(uint64(i.t.UnixMilli())) }

// Equals compares by epoch milliseconds. After UTC+millis normalisation
// in the constructors this matches t.Equal, but expressing it directly in
// terms of UnixMilli decouples value identity from time.Time internals
// (e.g. any monotonic clock reading).
func (i *Instant) Equals(other Value) bool {
	o, ok := other.(*Instant)
	return ok && i.t.UnixMilli() == o.t.UnixMilli()
}

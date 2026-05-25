package vm

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type theUUIDType struct{}

func (t *theUUIDType) String() string  { return t.Name() }
func (t *theUUIDType) Type() ValueType { return TypeType }
func (t *theUUIDType) Unbox() any      { return reflect.TypeFor[*theUUIDType]() }
func (t *theUUIDType) Name() string    { return "let-go.lang.UUID" }
func (t *theUUIDType) Box(bare any) (Value, error) {
	switch v := bare.(type) {
	case string:
		u := ParseUUID(v)
		if u == nil {
			return NIL, fmt.Errorf("invalid UUID: %s", v)
		}
		return u, nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", UUIDType)
}

// UUIDType is the type of UUID values
var UUIDType *theUUIDType = &theUUIDType{}

// UUID holds a UUID string in canonical lowercase form.
type UUID struct {
	val string // canonical lowercase 8-4-4-4-12 form
}

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ParseUUID parses a UUID string. Returns nil if invalid.
func ParseUUID(s string) *UUID {
	s = strings.ToLower(s)
	if !uuidRe.MatchString(s) {
		if normalized, ok := normalizeLenientUUID(s); ok {
			s = normalized
		}
	}
	if !uuidRe.MatchString(s) {
		return nil
	}
	return &UUID{val: s}
}

func normalizeLenientUUID(s string) (string, bool) {
	parts := strings.Split(s, "-")
	if len(parts) != 5 {
		return "", false
	}
	widths := []int{8, 4, 4, 4, 12}
	for i, p := range parts {
		if p == "" {
			return "", false
		}
		for _, r := range p {
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
				return "", false
			}
		}
		if len(p) > widths[i] && i != 3 {
			return "", false
		}
		if len(p) > widths[i] {
			p = p[len(p)-widths[i]:]
		}
		parts[i] = strings.Repeat("0", widths[i]-len(p)) + p
	}
	return strings.Join(parts, "-"), true
}

// NewUUID creates a UUID from an already-validated canonical string.
func NewUUID(s string) *UUID {
	return &UUID{val: s}
}

func (u *UUID) Type() ValueType { return UUIDType }
func (u *UUID) Unbox() any      { return u.val }
func (u *UUID) String() string  { return "#uuid \"" + u.val + "\"" }

// Hash implements Hashable.
func (u *UUID) Hash() uint32 { return hashString(u.val) }

func (u *UUID) Equals(other Value) bool {
	o, ok := other.(*UUID)
	return ok && u.val == o.val
}

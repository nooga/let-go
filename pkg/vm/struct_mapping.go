package vm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// fieldConverter converts a reflect.Value (a single struct field) to a VM Value
// without going through the generic BoxValue dispatch.
type fieldConverter func(reflect.Value) Value

// fieldDeconverter writes a VM Value back into a reflect.Value (struct field).
type fieldDeconverter func(reflect.Value, Value) error

// StructMapping holds the bidirectional mapping between a Go struct type and a RecordType.
type StructMapping struct {
	RecType    *RecordType
	GoType     reflect.Type
	fieldMap   []int              // recordFieldIdx → struct field index
	converters []fieldConverter   // recordFieldIdx → fast Go→Value converter
	deconvs    []fieldDeconverter // recordFieldIdx → fast Value→Go converter
}

var structMappings = map[reflect.Type]*StructMapping{}
var structMappingsByRecord = map[*RecordType]*StructMapping{}

// RegisterStructType registers a Go struct type and creates a corresponding RecordType.
// The name parameter sets the record type name. Field names are derived from exported
// struct fields: use `letgo:"name"` struct tag for custom names, otherwise fields are
// converted from CamelCase to kebab-case (e.g. FirstName → first-name).
// Fields tagged `letgo:"-"` are skipped.
func RegisterStructType(goType reflect.Type, name string) *StructMapping {
	if goType.Kind() == reflect.Pointer {
		goType = goType.Elem()
	}
	if goType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("RegisterStructType: expected struct, got %s", goType.Kind()))
	}

	if m, ok := structMappings[goType]; ok {
		return m
	}

	var keywords []Keyword
	var fieldIndices []int

	for i := 0; i < goType.NumField(); i++ {
		f := goType.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("letgo")
		if tag == "-" {
			continue
		}
		var kwName string
		if tag != "" {
			kwName = tag
		} else {
			kwName = camelToKebab(f.Name)
		}
		keywords = append(keywords, Keyword(kwName))
		fieldIndices = append(fieldIndices, i)
	}

	// Build per-field converters based on the Go type's Kind at registration time.
	converters := make([]fieldConverter, len(fieldIndices))
	deconvs := make([]fieldDeconverter, len(fieldIndices))
	for i, idx := range fieldIndices {
		f := goType.Field(idx)
		converters[i] = makeFieldConverter(f.Type)
		deconvs[i] = makeFieldDeconverter(f.Type)
	}

	rt := NewRecordType(name, keywords)
	m := &StructMapping{
		RecType:    rt,
		GoType:     goType,
		fieldMap:   fieldIndices,
		converters: converters,
		deconvs:    deconvs,
	}
	structMappings[goType] = m
	structMappingsByRecord[rt] = m
	return m
}

// makeFieldConverter returns a fast converter for a specific Go type.
// Falls back to BoxValue for types without a fast path.
func makeFieldConverter(t reflect.Type) fieldConverter {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(v reflect.Value) Value { return MakeInt(int(v.Int())) }
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(v reflect.Value) Value { return MakeInt(int(v.Uint())) }
	case reflect.Float32, reflect.Float64:
		return func(v reflect.Value) Value { return Float(v.Float()) }
	case reflect.Bool:
		return func(v reflect.Value) Value { return Boolean(v.Bool()) }
	case reflect.String:
		return func(v reflect.Value) Value { return String(v.String()) }
	case reflect.Interface:
		// If the field is already a Value (e.g. vm.Value), return it directly.
		valueInterface := reflect.TypeFor[Value]()
		if t.Implements(valueInterface) || t == valueInterface {
			return func(v reflect.Value) Value {
				if v.IsNil() {
					return NIL
				}
				return v.Interface().(Value)
			}
		}
		return func(v reflect.Value) Value {
			if v.IsNil() {
				return NIL
			}
			val, err := BoxValue(v.Elem())
			if err != nil {
				return NIL
			}
			return val
		}
	default:
		// Fallback: use the generic path
		return func(v reflect.Value) Value {
			val, err := BoxValue(v)
			if err != nil {
				return NIL
			}
			return val
		}
	}
}

// makeFieldDeconverter returns a fast deconverter for a specific Go type.
func makeFieldDeconverter(t reflect.Type) fieldDeconverter {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(target reflect.Value, val Value) error {
			if v, ok := val.(Int); ok {
				target.SetInt(int64(v))
				return nil
			}
			if v, ok := val.(Float); ok {
				target.SetInt(int64(v))
				return nil
			}
			return fmt.Errorf("expected Int, got %s", val.Type().Name())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(target reflect.Value, val Value) error {
			if v, ok := val.(Int); ok {
				target.SetUint(uint64(v))
				return nil
			}
			return fmt.Errorf("expected Int, got %s", val.Type().Name())
		}
	case reflect.Float32, reflect.Float64:
		return func(target reflect.Value, val Value) error {
			if v, ok := val.(Float); ok {
				target.SetFloat(float64(v))
				return nil
			}
			if v, ok := val.(Int); ok {
				target.SetFloat(float64(v))
				return nil
			}
			return fmt.Errorf("expected Float, got %s", val.Type().Name())
		}
	case reflect.Bool:
		return func(target reflect.Value, val Value) error {
			if v, ok := val.(Boolean); ok {
				target.SetBool(bool(v))
				return nil
			}
			return fmt.Errorf("expected Boolean, got %s", val.Type().Name())
		}
	case reflect.String:
		return func(target reflect.Value, val Value) error {
			if v, ok := val.(String); ok {
				target.SetString(string(v))
				return nil
			}
			if v, ok := val.(Keyword); ok {
				target.SetString(string(v))
				return nil
			}
			return fmt.Errorf("expected String, got %s", val.Type().Name())
		}
	case reflect.Interface:
		valueInterface := reflect.TypeFor[Value]()
		if t.Implements(valueInterface) || t == valueInterface {
			return func(target reflect.Value, val Value) error {
				target.Set(reflect.ValueOf(val))
				return nil
			}
		}
		return func(target reflect.Value, val Value) error {
			return unboxInto(target, val)
		}
	default:
		// Fallback to the generic unboxInto
		return func(target reflect.Value, val Value) error {
			return unboxInto(target, val)
		}
	}
}

// StructToRecord converts a Go struct value to a Record.
// The original value is stored for fast roundtrip back to Go.
// Uses cached per-field converters to avoid BoxValue dispatch overhead.
func (m *StructMapping) StructToRecord(v any) *Record {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	fields := make([]Value, len(m.RecType.fields))
	for i, structIdx := range m.fieldMap {
		fields[i] = m.converters[i](rv.Field(structIdx))
	}

	return &Record{
		rtype:  m.RecType,
		fields: fields,
		extra:  EmptyPersistentMap,
		origin: v,
	}
}

// RecordToStruct populates a Go struct from a Record's fields.
// If the Record has an origin of the same type, returns it directly (fast path).
func (m *StructMapping) RecordToStruct(r *Record, target any) error {
	// Fast path: if record has origin of same type, copy it
	if r.origin != nil {
		ov := reflect.ValueOf(r.origin)
		if ov.Kind() == reflect.Pointer {
			ov = ov.Elem()
		}
		if ov.Type() == m.GoType {
			tv := reflect.ValueOf(target)
			if tv.Kind() != reflect.Pointer {
				return fmt.Errorf("target must be a pointer")
			}
			tv.Elem().Set(ov)
			return nil
		}
	}

	// Slow path: read fields from the Record using cached deconverters
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Pointer {
		return fmt.Errorf("target must be a pointer")
	}
	tv = tv.Elem()

	for i, structIdx := range m.fieldMap {
		val := r.fields[i]
		if val == nil || val == NIL {
			continue
		}
		if err := m.deconvs[i](tv.Field(structIdx), val); err != nil {
			return fmt.Errorf("field %s: %w", m.RecType.fields[i], err)
		}
	}
	return nil
}

// LookupStructMapping returns the mapping for a Go type, or nil if not registered.
func LookupStructMapping(goType reflect.Type) *StructMapping {
	if goType.Kind() == reflect.Pointer {
		goType = goType.Elem()
	}
	return structMappings[goType]
}

// LookupStructMappingByRecord returns the mapping for a RecordType, or nil.
func LookupStructMappingByRecord(rt *RecordType) *StructMapping {
	return structMappingsByRecord[rt]
}

// unboxInto sets a reflect.Value from a VM Value using type-appropriate conversion.
func unboxInto(target reflect.Value, val Value) error {
	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, ok := val.(Int); ok {
			target.SetInt(int64(v))
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, ok := val.(Int); ok {
			target.SetUint(uint64(v))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		switch v := val.(type) {
		case Float:
			target.SetFloat(float64(v))
			return nil
		case Int:
			target.SetFloat(float64(v))
			return nil
		}
	case reflect.Bool:
		if v, ok := val.(Boolean); ok {
			target.SetBool(bool(v))
			return nil
		}
	case reflect.String:
		switch v := val.(type) {
		case String:
			target.SetString(string(v))
			return nil
		case Keyword:
			target.SetString(string(v))
			return nil
		}
	case reflect.Slice:
		if sq, ok := val.(Sequable); ok {
			return unboxSliceInto(target, sq.Seq())
		}
	case reflect.Map:
		// Guarded by the same explicit type switch unboxMapInto uses, so a
		// value that is not a let-go map (a *Boxed already holding a Go map,
		// say) falls through to the assignability fallback below.
		switch val.(type) {
		case Map, *PersistentMap, *SortedMap:
			return unboxMapInto(target, val)
		}
	case reflect.Pointer:
		// If it's a Boxed value holding a pointer of the right type, unwrap it
		if b, ok := val.(*Boxed); ok {
			bv := reflect.ValueOf(b.Unbox())
			if bv.Type().AssignableTo(target.Type()) {
				target.Set(bv)
				return nil
			}
		}
	case reflect.Interface:
		// `any` holds every Go value, Go nil included, so converting into it
		// is TOTAL — it must never reach the error return below. When it did,
		// the caller (boxArgForReflect) treated the failure as "this sequence
		// has no Go form", handed the whole slice over as []vm.Value, and the
		// reflect call died with:
		//
		//	reflect: Call using []vm.Value as type []interface {}
		//
		// A let-go nil was the common trigger, which made every SQL NULL
		// parameter fail. Restricted to the empty interface: a non-empty one
		// (error, io.Reader, …) genuinely may not be satisfiable, and the
		// paths below still decide that.
		if target.NumMethod() == 0 {
			if raw := val.Unbox(); raw != nil {
				target.Set(reflect.ValueOf(raw))
			} else {
				// Zero explicitly rather than leaving the target untouched.
				// unboxSliceInto allocates a fresh element, but unboxInto is
				// shared — RecordToStruct writes into an EXISTING struct
				// field, which may already hold a value. Skipping the write
				// would report success while leaving stale data behind. More
				// than NIL lands here: any value whose Unbox yields nil.
				target.Set(reflect.Zero(target.Type()))
			}
			return nil
		}
	}

	// Fallback: try Unbox and assignability
	raw := val.Unbox()
	rv := reflect.ValueOf(raw)
	if rv.IsValid() && rv.Type().AssignableTo(target.Type()) {
		target.Set(rv)
		return nil
	}
	if rv.IsValid() && rv.Type().ConvertibleTo(target.Type()) {
		target.Set(rv.Convert(target.Type()))
		return nil
	}

	return fmt.Errorf("cannot convert %s (%s) to %s", val, val.Type().Name(), target.Type())
}

func unboxSliceInto(target reflect.Value, s Seq) error {
	elemType := target.Type().Elem()
	var result []reflect.Value
	for !SeqIsEmpty(s) {
		elem := reflect.New(elemType).Elem()
		if err := unboxInto(elem, s.First()); err != nil {
			return err
		}
		result = append(result, elem)
		s = s.Next()
	}
	slice := reflect.MakeSlice(target.Type(), len(result), len(result))
	for i, v := range result {
		slice.Index(i).Set(v)
	}
	target.Set(slice)
	return nil
}


// unboxMapInto converts a let-go map into a Go map target, key by key and value
// by value, so map[string]any and map[string]string both work and a map nested
// in a struct field or slice element converts too.
//
// The source is identified by an explicit TYPE switch rather than by inspecting
// the first element for a MapEntry: every let-go map returns EmptyList from
// Seq() when empty, so an empty map and an empty vector are indistinguishable
// that way. TransientMap is deliberately absent — it is a builder, not a value
// handed to Go.
//
// Keyword keys become string keys because unboxInto's reflect.String case
// already accepts a Keyword and a Keyword stores its name without the leading
// colon. That follows an existing convention rather than inventing one.
func unboxMapInto(target reflect.Value, val Value) error {
	var s Seq
	switch m := val.(type) {
	case Map:
		s = m.Seq()
	case *PersistentMap:
		s = m.Seq()
	case *SortedMap:
		s = m.Seq()
	default:
		return fmt.Errorf("cannot convert %s (%s) to %s", val, val.Type().Name(), target.Type())
	}

	mapType := target.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()
	out := reflect.MakeMap(mapType)
	for !SeqIsEmpty(s) {
		entry, ok := s.First().(MapEntry)
		if !ok {
			return fmt.Errorf("cannot convert %s to %s: %s is not a map entry", val.Type().Name(), mapType, s.First())
		}
		if keyType.Kind() == reflect.String && unboxesToInteger(entry.Key) {
			// Go reports int64 as ConvertibleTo string and converts it to a
			// RUNE, so unboxInto's generic fallback would silently turn the
			// let-go key 65 into the Go key "A". Reject instead — a wrong key
			// that looks plausible is worse than a failure — and keep the
			// rejection local to the map path so no shared conversion
			// behaviour changes.
			return fmt.Errorf("cannot convert map key %s (%s) to %s", entry.Key, entry.Key.Type().Name(), keyType)
		}
		k := reflect.New(keyType).Elem()
		if err := unboxInto(k, entry.Key); err != nil {
			return err
		}
		v := reflect.New(elemType).Elem()
		if err := unboxInto(v, entry.Value); err != nil {
			return err
		}
		if err := setMapKey(out, k, v); err != nil {
			return fmt.Errorf("cannot use %s (%s) as a key in %s: %w",
				entry.Key, entry.Key.Type().Name(), mapType, err)
		}
		s = s.Next()
	}
	target.Set(out)
	return nil
}

// setMapKey stores one converted entry, turning an unhashable key into a
// conversion error rather than a panic.
//
// An interface-kind key (map[any]any) accepts any dynamic type, and unboxInto
// puts a []vm.Value there for a let-go vector key — SetMapIndex then panics
// with "hash of unhashable type". That has to be caught here: unboxMapInto is
// reached from RecordToStruct too, which has no recover of its own, so a panic
// is not guaranteed to be contained.
//
// The check is on the VALUE, not the type, because reflect.Type.Comparable is
// not enough: a struct{ X any } holding a slice reports Comparable() == true
// and still panics when hashed. Only an attempted insert settles it.
func setMapKey(out, k, v reflect.Value) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("the value is not hashable")
		}
	}()
	out.SetMapIndex(k, v)
	return nil
}

// unboxesToInteger reports whether a value reaches Go as an integer, and so
// would be rune-converted on the way into a string-kind target.
func unboxesToInteger(v Value) bool {
	raw := v.Unbox()
	if raw == nil {
		return false
	}
	switch reflect.ValueOf(raw).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

// camelToKebab converts CamelCase to kebab-case.
// e.g. "FirstName" → "first-name", "HTTPServer" → "http-server", "ID" → "id"
func camelToKebab(s string) string {
	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				// Don't add hyphen if previous was also upper and next is upper or end
				// (handles "HTTP" → "http" not "h-t-t-p")
				prevUpper := unicode.IsUpper(runes[i-1])
				nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
				if !prevUpper || nextLower {
					result.WriteRune('-')
				}
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

package vm

import (
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
	if goType.Kind() == reflect.Ptr {
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
func (m *StructMapping) StructToRecord(v interface{}) *Record {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
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
func (m *StructMapping) RecordToStruct(r *Record, target interface{}) error {
	// Fast path: if record has origin of same type, copy it
	if r.origin != nil {
		ov := reflect.ValueOf(r.origin)
		if ov.Kind() == reflect.Ptr {
			ov = ov.Elem()
		}
		if ov.Type() == m.GoType {
			tv := reflect.ValueOf(target)
			if tv.Kind() != reflect.Ptr {
				return fmt.Errorf("target must be a pointer")
			}
			tv.Elem().Set(ov)
			return nil
		}
	}

	// Slow path: read fields from the Record using cached deconverters
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Ptr {
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
	if goType.Kind() == reflect.Ptr {
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
	case reflect.Ptr:
		// If it's a Boxed value holding a pointer of the right type, unwrap it
		if b, ok := val.(*Boxed); ok {
			bv := reflect.ValueOf(b.Unbox())
			if bv.Type().AssignableTo(target.Type()) {
				target.Set(bv)
				return nil
			}
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
	for s != nil {
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

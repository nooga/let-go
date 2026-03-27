package vm

import "reflect"

// RegisterStruct registers a Go struct type T and returns its StructMapping.
// Usage: mapping := vm.RegisterStruct[MyStruct]("my-ns/MyStruct")
func RegisterStruct[T any](name string) *StructMapping {
	var zero T
	return RegisterStructType(reflect.TypeOf(zero), name)
}

// ToRecord converts a Go struct to a Record using its registered mapping.
// Panics if the struct type is not registered.
func ToRecord[T any](v T) *Record {
	m := LookupStructMapping(reflect.TypeOf(v))
	if m == nil {
		panic("struct type not registered: " + reflect.TypeOf(v).String())
	}
	return m.StructToRecord(v)
}

// ToStruct converts a Record back to a Go struct T.
// Uses the fast path (returns original) if the Record hasn't been mutated.
func ToStruct[T any](r *Record) (T, error) {
	var result T
	m := LookupStructMapping(reflect.TypeOf(result))
	if m == nil {
		return result, NewExecutionError("struct type not registered: " + reflect.TypeOf(result).String())
	}
	err := m.RecordToStruct(r, &result)
	return result, err
}

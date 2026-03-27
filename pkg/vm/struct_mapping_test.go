package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelToKebab(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Name", "name"},
		{"FirstName", "first-name"},
		{"ID", "id"},
		{"HTTPServer", "http-server"},
		{"userID", "user-id"},
		{"X", "x"},
		{"XY", "xy"},
		{"getHTTPResponse", "get-http-response"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, camelToKebab(tt.input))
		})
	}
}

type TestPoint struct {
	X int
	Y int
}

type TestPerson struct {
	FirstName string `letgo:"first-name"`
	LastName  string `letgo:"last-name"`
	Age       int
}

type TestWithSkip struct {
	Name    string
	Secret  string `letgo:"-"`
	Visible int
}

func TestRegisterStruct(t *testing.T) {
	m := RegisterStruct[TestPoint]("test/Point")
	assert.Equal(t, "test/Point", m.RecType.Name())
	assert.Equal(t, []Keyword{"x", "y"}, m.RecType.Fields())

	// Same type returns same mapping
	m2 := RegisterStruct[TestPoint]("test/Point")
	assert.Same(t, m, m2)
}

func TestRegisterStructWithTags(t *testing.T) {
	m := RegisterStruct[TestPerson]("test/Person")
	assert.Equal(t, []Keyword{"first-name", "last-name", "age"}, m.RecType.Fields())
}

func TestRegisterStructSkipField(t *testing.T) {
	m := RegisterStruct[TestWithSkip]("test/WithSkip")
	assert.Equal(t, []Keyword{"name", "visible"}, m.RecType.Fields())
}

func TestStructToRecord(t *testing.T) {
	m := RegisterStruct[TestPoint]("test/Point")
	p := TestPoint{X: 3, Y: 4}
	r := m.StructToRecord(p)

	assert.Equal(t, m.RecType, r.Type())
	assert.Equal(t, Int(3), r.ValueAt(Keyword("x")))
	assert.Equal(t, Int(4), r.ValueAt(Keyword("y")))

	// Origin is preserved
	assert.NotNil(t, r.Origin())
}

func TestStructToRecordPointer(t *testing.T) {
	m := RegisterStruct[TestPerson]("test/Person")
	p := &TestPerson{FirstName: "Alice", LastName: "Smith", Age: 30}
	r := m.StructToRecord(p)

	assert.Equal(t, String("Alice"), r.ValueAt(Keyword("first-name")))
	assert.Equal(t, String("Smith"), r.ValueAt(Keyword("last-name")))
	assert.Equal(t, Int(30), r.ValueAt(Keyword("age")))
}

func TestRecordToStructFastPath(t *testing.T) {
	RegisterStruct[TestPoint]("test/Point")
	original := TestPoint{X: 10, Y: 20}
	r := ToRecord(original)

	// Fast path: origin present, no mutation
	result, err := ToStruct[TestPoint](r)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestRecordToStructSlowPath(t *testing.T) {
	RegisterStruct[TestPoint]("test/Point")
	original := TestPoint{X: 10, Y: 20}
	r := ToRecord(original)

	// Mutate via assoc — clears origin
	mutated := r.Assoc(Keyword("x"), Int(99)).(*Record)
	assert.Nil(t, mutated.Origin())

	result, err := ToStruct[TestPoint](mutated)
	assert.NoError(t, err)
	assert.Equal(t, TestPoint{X: 99, Y: 20}, result)
}

func TestRecordToStructWithStrings(t *testing.T) {
	RegisterStruct[TestPerson]("test/Person")
	r := ToRecord(TestPerson{FirstName: "Bob", LastName: "Jones", Age: 25})

	// Mutate last name
	mutated := r.Assoc(Keyword("last-name"), String("Smith")).(*Record)

	result, err := ToStruct[TestPerson](mutated)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", result.FirstName)
	assert.Equal(t, "Smith", result.LastName)
	assert.Equal(t, 25, result.Age)
}

func TestRoundtripPreservesIdentity(t *testing.T) {
	RegisterStruct[TestPoint]("test/Point")
	original := TestPoint{X: 42, Y: 99}

	// Go → Record → Go should be exact roundtrip
	r := ToRecord(original)
	result, err := ToStruct[TestPoint](r)
	assert.NoError(t, err)
	assert.Equal(t, original, result)

	// Unbox should return the original too
	assert.Equal(t, original, r.Unbox())
}

func TestRecordFromStructWorksAsMap(t *testing.T) {
	RegisterStruct[TestPoint]("test/Point")
	r := ToRecord(TestPoint{X: 1, Y: 2})

	// Keyword access works
	assert.Equal(t, Int(1), r.ValueAt(Keyword("x")))

	// Record-as-function works
	v, err := r.Invoke([]Value{Keyword("y")})
	assert.NoError(t, err)
	assert.Equal(t, Int(2), v)

	// Interop .field works
	v, err = r.InvokeMethod("x", nil)
	assert.NoError(t, err)
	assert.Equal(t, Int(1), v)

	// Count works
	assert.Equal(t, 2, r.RawCount())

	// Seq works
	s := r.Seq()
	assert.NotNil(t, s)
}

func TestBoxValueAutoConvertsRegisteredStruct(t *testing.T) {
	RegisterStruct[TestPoint]("test/Point")
	p := TestPoint{X: 7, Y: 8}

	val, err := ToLetGo(p)
	assert.NoError(t, err)

	r, ok := val.(*Record)
	assert.True(t, ok, "expected Record, got %T", val)
	assert.Equal(t, Int(7), r.ValueAt(Keyword("x")))
	assert.Equal(t, Int(8), r.ValueAt(Keyword("y")))
}

type TestNested struct {
	Name  string
	Point TestPoint
}

func TestNestedStructConversion(t *testing.T) {
	RegisterStruct[TestPoint]("test/Point")
	RegisterStruct[TestNested]("test/Nested")

	n := TestNested{Name: "origin", Point: TestPoint{X: 1, Y: 2}}
	r := ToRecord(n)

	assert.Equal(t, String("origin"), r.ValueAt(Keyword("name")))
	inner, ok := r.ValueAt(Keyword("point")).(*Record)
	assert.True(t, ok, "nested struct should be a Record")
	assert.Equal(t, Int(1), inner.ValueAt(Keyword("x")))
	assert.Equal(t, Int(2), inner.ValueAt(Keyword("y")))
}

type TestWithFloat struct {
	Value float64
	Label string
}

func TestStructWithFloat(t *testing.T) {
	RegisterStruct[TestWithFloat]("test/WithFloat")
	original := TestWithFloat{Value: 3.14, Label: "pi"}

	r := ToRecord(original)
	assert.Equal(t, Float(3.14), r.ValueAt(Keyword("value")))
	assert.Equal(t, String("pi"), r.ValueAt(Keyword("label")))

	result, err := ToStruct[TestWithFloat](r)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

package vm

import "testing"

func TestMetaValueDecoratesAnyValueType(t *testing.T) {
	meta := NewArrayMap([]Value{Keyword("source"), String("test")})
	wrapped := NewMetaValue(String("value"), meta)

	if wrapped.Type() != StringType || wrapped.String() != String("value").String() || wrapped.Unbox() != "value" {
		t.Fatalf("generic metadata decorator changed wrapped value behavior: %v", wrapped)
	}
	if wrapped.Meta() != meta {
		t.Fatalf("metadata = %v, want %v", wrapped.Meta(), meta)
	}

	replacement := NewArrayMap([]Value{Keyword("source"), String("replacement")})
	rewrapped := wrapped.WithMeta(replacement).(*MetaValue[String])
	if rewrapped.Meta() != replacement || rewrapped.Wrapped() != String("value") {
		t.Fatalf("WithMeta did not preserve value and replace metadata: %v", rewrapped)
	}
}

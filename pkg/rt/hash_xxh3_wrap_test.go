package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/xxh3"
)

// The Wrap-based adapters replace the generated bindings under TinyGo wasm, so
// they must return the library's digest boxed the way IntType.Box boxes a
// uint64 (reinterpreted as int64), and reject the wrong arity like a
// generated binding does.
func TestXxh3WrapMatchesLibraryAndBoxing(t *testing.T) {
	inputs := []string{"", "a", "hello", "The quick brown fox jumps over the lazy dog", string(make([]byte, 300))}
	seeds := []uint64{0, 1, 12345, 1 << 40, 1<<64 - 1}
	for _, in := range inputs {
		for _, seed := range seeds {
			want := int64(xxh3.HashStringSeed(in, seed))
			got, err := xxh3HashStringSeed([]vm.Value{vm.String(in), vm.Int(int64(seed))})
			if err != nil {
				t.Fatalf("HashStringSeed(%q, %d): %v", in, seed, err)
			}
			if int64(got.(vm.Int)) != want {
				t.Fatalf("HashStringSeed(%q, %d) = %d, want %d", in, seed, got, want)
			}
			gotB, err := xxh3HashSeed([]vm.Value{vm.NewByteArrayFrom([]byte(in)), vm.Int(int64(seed))})
			if err != nil {
				t.Fatalf("HashSeed(%q, %d): %v", in, seed, err)
			}
			if gotB != got {
				t.Fatalf("HashSeed and HashStringSeed disagree on %q/%d: %v vs %v", in, seed, gotB, got)
			}
		}
		if got, _ := xxh3HashString([]vm.Value{vm.String(in)}); int64(got.(vm.Int)) != int64(xxh3.HashString(in)) {
			t.Fatalf("HashString(%q) = %v", in, got)
		}
		if got, _ := xxh3Hash([]vm.Value{vm.NewByteArrayFrom([]byte(in))}); int64(got.(vm.Int)) != int64(xxh3.Hash([]byte(in))) {
			t.Fatalf("Hash(%q) = %v", in, got)
		}
	}
	if _, err := xxh3HashSeed([]vm.Value{vm.Int(1), vm.Int(0)}); err == nil {
		t.Fatal("expected a type error for a non-bytes input")
	}
	if _, err := xxh3HashString([]vm.Value{vm.String("a"), vm.Keyword("ignored")}); err == nil {
		t.Fatal("expected an arity error for an extra argument")
	}
	boxed, err := vm.IntType.Box(xxh3.HashString("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := xxh3HashString([]vm.Value{vm.String("hello")}); got != boxed {
		t.Fatalf("adapter boxes %v, IntType.Box boxes %v", got, boxed)
	}
}

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// base64-encode/decode are rt builtins, so exercise them through the compiler the
// way reader_bb_test does for other rt features.
func TestBase64(t *testing.T) {
	// byte-array → padded base64 String (PNG signature → "iVBORw==").
	if v := evalStr(t, `(base64-encode (byte-array [137 80 78 71]))`); v != vm.String("iVBORw==") {
		t.Fatalf("encode byte-array: got %v", v)
	}
	// String input encodes too.
	if v := evalStr(t, `(base64-encode "hi")`); v != vm.String("aGk=") {
		t.Fatalf("encode string: got %v", v)
	}
	// Round-trip: decode(encode(bytes)) preserves the byte count.
	if v := evalStr(t, `(count (base64-decode (base64-encode (byte-array (range 100)))))`); v != vm.Int(100) {
		t.Fatalf("round-trip count: got %v", v)
	}
}

func TestBase64URL(t *testing.T) {
	// URL-safe alphabet: 0xFF*3 fills all six-bit groups with index 63, which is
	// '_' in base64url (std base64 would also use '/' here — this pins -_ over +/).
	if v := evalStr(t, `(base64url-encode (byte-array [255 255 255]))`); v != vm.String("____") {
		t.Fatalf("encode url-safe alphabet: got %v", v)
	}
	// No padding, and index 62 is '-' (std base64 emits "+A==").
	if v := evalStr(t, `(base64url-encode (byte-array [248]))`); v != vm.String("-A") {
		t.Fatalf("encode no-pad: got %v", v)
	}
	// Round-trip: decode(encode(bytes)) preserves the byte count.
	if v := evalStr(t, `(count (base64url-decode (base64url-encode (byte-array (range 100)))))`); v != vm.Int(100) {
		t.Fatalf("round-trip count: got %v", v)
	}
}

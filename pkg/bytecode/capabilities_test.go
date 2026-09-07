package bytecode

import (
	"fmt"
	"strings"
	"testing"
)

func TestDescribeCapabilities(t *testing.T) {
	got := DescribeCapabilities(0)
	if got != "(none)" {
		t.Fatalf("empty mask: got %q", got)
	}
	got = DescribeCapabilities(CapOpcodeSet)
	if !strings.Contains(got, "CapOpcodeSet") || !strings.Contains(got, "since v1.12.0") {
		t.Fatalf("named cap: got %q", got)
	}
	got = DescribeCapabilities(0x2)
	if !strings.Contains(got, "unknown bit 1") {
		t.Fatalf("unknown bit: got %q", got)
	}
}

func TestUnsupportedCapabilityErrorNamedMinVersion(t *testing.T) {
	// CapOpcodeSet is supported today; synthesize a "named but unsupported"
	// case by describing a mask that includes only CapOpcodeSet against a
	// zero support set via MinVersionForCapabilities.
	if min := MinVersionForCapabilities(CapOpcodeSet); min != "1.12.0" {
		t.Fatalf("min version: got %q, want 1.12.0", min)
	}
	if min := MinVersionForCapabilities(0x2); min != "" {
		t.Fatalf("unknown-only mask should have no min version, got %q", min)
	}
}

func TestUnsupportedCapabilityErrorMaskIsUnsupportedSubset(t *testing.T) {
	// A bundle asking for a supported bit (CapOpcodeSet) AND an unknown bit
	// must report only the unsupported subset in the "unsupported" mask —
	// not the full requested mask.
	err := UnsupportedCapabilityError(CapOpcodeSet | 0x2)
	msg := err.Error()
	if !strings.Contains(msg, "unsupported capabilities: unknown bit 1 (0x2) (mask 0x00000002)") {
		t.Fatalf("unsupported section should describe only bit 1 with mask 0x2, got: %s", msg)
	}
	if strings.Contains(msg, "mask 0x00000003") {
		t.Fatalf("unsupported mask must not include the supported bit, got: %s", msg)
	}
}

func TestFormatVersionReport(t *testing.T) {
	got := FormatVersionReport("lg", "1.99.0 (deadbee)")
	for _, want := range []string{
		"lg 1.99.0 (deadbee)\n",
		fmt.Sprintf("lgb: format %d (default write), %d (max)\n", uncompressedFormatVersion, FormatVersion),
		"capabilities: CapOpcodeSet",
		"opcodes:",
		"signature",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}

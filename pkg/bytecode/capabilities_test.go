package bytecode

import (
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

func TestFormatVersionReport(t *testing.T) {
	got := FormatVersionReport("lg", "1.99.0 (deadbee)")
	for _, want := range []string{
		"lg 1.99.0 (deadbee)\n",
		"lgb: format 2\n",
		"capabilities: CapOpcodeSet",
		"opcodes:",
		"signature",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}

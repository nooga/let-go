package bytecode

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
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
	got = DescribeCapabilities(1 << 31)
	if !strings.Contains(got, "unknown bit 31") {
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
	if min := MinVersionForCapabilities(1 << 31); min != "" {
		t.Fatalf("unknown-only mask should have no min version, got %q", min)
	}
}

func TestUnsupportedCapabilityErrorMaskIsUnsupportedSubset(t *testing.T) {
	// A bundle asking for a supported bit (CapOpcodeSet) AND an unknown bit
	// must report only the unsupported subset in the "unsupported" mask —
	// not the full requested mask.
	err := UnsupportedCapabilityError(CapOpcodeSet | 1<<31)
	msg := err.Error()
	if !strings.Contains(msg, "unsupported capabilities: unknown bit 31 (0x80000000) (mask 0x80000000)") {
		t.Fatalf("unsupported section should describe only bit 31 with mask 0x80000000, got: %s", msg)
	}
	if strings.Contains(msg, "mask 0x80000001") {
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

// moduleCapabilities builds a module holding consts, round-trips it through
// Encode and Decode, and returns the decoded capability mask.
func moduleCapabilities(t *testing.T, consts ...vm.Value) uint32 {
	t.Helper()
	b := NewModuleBuilder()
	for _, c := range consts {
		b.AddConst(c)
	}
	var buf bytes.Buffer
	if err := Encode(&buf, b.Build()); err != nil {
		t.Fatal(err)
	}
	m, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	return m.Capabilities
}

func TestTypedConstantsCapabilityOnlyWhenNeeded(t *testing.T) {
	pt := vm.DefineRecordType("captest.ns", "Pt", []vm.Keyword{"x"})
	box := vm.DefineDType("captest.ns", "Box", []vm.Symbol{"v"})
	bare := vm.NewRecordType("Pt", []vm.Keyword{"x"})
	for _, tc := range []struct {
		name   string
		consts []vm.Value
		want   bool
	}{
		{"plain constants", []vm.Value{vm.Int(1), vm.String("s")}, false},
		{"record of a type with no namespace", []vm.Value{vm.NewRecord(bare, vm.EmptyPersistentMap)}, false},
		{"defrecord instance", []vm.Value{vm.NewRecord(pt, vm.EmptyPersistentMap)}, true},
		{"defrecord type", []vm.Value{pt}, true},
		{"deftype instance", []vm.Value{vm.NewDTypeInstance(box, []vm.Value{vm.Int(7)})}, true},
		{"defrecord instance in a vector", []vm.Value{vm.NewArrayVector([]vm.Value{vm.NewRecord(pt, vm.EmptyPersistentMap)})}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			caps := moduleCapabilities(t, tc.consts...)
			if got := caps&CapTypedConstants != 0; got != tc.want {
				t.Fatalf("CapTypedConstants set = %v, want %v (mask 0x%x)", got, tc.want, caps)
			}
			if caps&CapOpcodeSet == 0 {
				t.Fatalf("CapOpcodeSet dropped (mask 0x%x)", caps)
			}
		})
	}
}

// A reify instance cannot be written, so it must not set the bit either.
func TestTypedConstantsCapabilityNotSetForReify(t *testing.T) {
	b := NewModuleBuilder()
	b.AddConst(vm.NewDTypeInstance(vm.NewDType("reify__1", nil), nil))
	if caps := b.Build().Capabilities; caps&CapTypedConstants != 0 {
		t.Fatalf("CapTypedConstants set for a reify instance (mask 0x%x)", caps)
	}
}

func TestTypedConstantsCapabilityIsNamed(t *testing.T) {
	if got := DescribeCapabilities(CapTypedConstants); got != "CapTypedConstants (0x2, since v1.14.0)" {
		t.Fatalf("describe: got %q", got)
	}
	if min := MinVersionForCapabilities(CapOpcodeSet | CapTypedConstants); min != "1.14.0" {
		t.Fatalf("min version: got %q, want 1.14.0", min)
	}
	if SupportedCapabilities&CapTypedConstants == 0 {
		t.Fatal("this runtime does not accept CapTypedConstants")
	}
	if got := FormatVersionReport("lg", "x"); !strings.Contains(got, "CapTypedConstants") {
		t.Fatalf("lg -v does not list CapTypedConstants:\n%s", got)
	}
}

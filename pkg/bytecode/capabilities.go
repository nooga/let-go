package bytecode

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// CapabilityInfo names a known capability bit and the lg release that first
// accepted it. Keep this in sync when adding bits to the Capability constants
// in tags.go — the reject path and `lg -v` both read from here.
type CapabilityInfo struct {
	Bit   uint32
	Name  string
	Since string // semver without leading "v", e.g. "1.12.0"
}

// KnownCapabilities lists every capability bit this tree knows about, in bit
// order. SupportedCapabilities is the subset this runtime accepts.
var KnownCapabilities = []CapabilityInfo{
	{Bit: CapOpcodeSet, Name: "CapOpcodeSet", Since: "1.12.0"},
}

// SupportedCapabilities is the capability mask this runtime accepts.
// Must match the check in readCapabilities.
const SupportedCapabilities = CapOpcodeSet

// DescribeCapabilities formats a capability mask for humans: named bits as
// CapName (0x…, since vX.Y.Z), unknown bits as "unknown bit N (0x…)".
func DescribeCapabilities(mask uint32) string {
	if mask == 0 {
		return "(none)"
	}
	var parts []string
	seen := uint32(0)
	for _, info := range KnownCapabilities {
		if mask&info.Bit == 0 {
			continue
		}
		seen |= info.Bit
		parts = append(parts, fmt.Sprintf("%s (0x%x, since v%s)", info.Name, info.Bit, info.Since))
	}
	unknown := mask &^ seen
	for bit := 0; bit < 32 && unknown != 0; bit++ {
		b := uint32(1) << uint(bit)
		if unknown&b == 0 {
			continue
		}
		parts = append(parts, fmt.Sprintf("unknown bit %d (0x%x)", bit, b))
		unknown &^= b
	}
	return strings.Join(parts, ", ")
}

// MinVersionForCapabilities returns the highest Since among named bits in
// mask, or "" when mask has no named capabilities (unknown bits alone).
func MinVersionForCapabilities(mask uint32) string {
	var min string
	for _, info := range KnownCapabilities {
		if mask&info.Bit == 0 {
			continue
		}
		if min == "" || versionLess(min, info.Since) {
			min = info.Since
		}
	}
	return min
}

// versionLess reports whether a < b for dotted numeric semver (no pre-release).
func versionLess(a, b string) bool {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		ai, bi := 0, 0
		if i < len(as) {
			ai, _ = strconv.Atoi(as[i])
		}
		if i < len(bs) {
			bi, _ = strconv.Atoi(bs[i])
		}
		if ai != bi {
			return ai < bi
		}
	}
	return false
}

// UnsupportedCapabilityError builds the user-facing reject for a capability
// mask that asks for bits this runtime does not support.
func UnsupportedCapabilityError(caps uint32) error {
	unsupported := caps &^ SupportedCapabilities
	msg := fmt.Sprintf(
		"unsupported capabilities: %s (mask 0x%08x); runtime supports: %s (mask 0x%08x)",
		DescribeCapabilities(unsupported), unsupported,
		DescribeCapabilities(SupportedCapabilities), SupportedCapabilities)
	if min := MinVersionForCapabilities(unsupported); min != "" {
		msg += fmt.Sprintf("; needs lg >= v%s", min)
	}
	msg += " — recompile the bundle with a matching lg or run it on a newer runtime"
	return errors.New(msg)
}

// FormatVersionReport returns the multi-line -v/--version body for tool
// (e.g. "lg" or "lg-runtime") at the given version string.
//
// The lgb line reports both the default write version (plain bundles stay on
// uncompressedFormatVersion) and FormatVersion (the newest this tree reads and
// writes when a newer-only flag such as FlagCompressed is set).
func FormatVersionReport(tool, versionStr string) string {
	count, hash := vm.OpcodeSetSignature()
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s\n", tool, versionStr)
	fmt.Fprintf(&b, "lgb: format %d (default write), %d (max)\n", uncompressedFormatVersion, FormatVersion)
	fmt.Fprintf(&b, "capabilities: %s\n", DescribeCapabilities(SupportedCapabilities))
	fmt.Fprintf(&b, "opcodes: %d (signature %016x)\n", count, hash)
	return b.String()
}

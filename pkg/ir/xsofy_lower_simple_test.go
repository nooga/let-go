//go:build bootstrap

package ir_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// ExtractDefns extracts all top-level (defn ...) forms from Lisp source
func ExtractDefnsSimple(source string) []string {
	var defns []string
	lines := strings.Split(source, "\n")
	var currentDefn strings.Builder
	parenDepth := 0
	inDefn := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of a new defn
		if !inDefn && strings.HasPrefix(trimmed, "(defn ") {
			inDefn = true
			parenDepth = 0
		}

		if inDefn {
			currentDefn.WriteString(line)
			currentDefn.WriteString("\n")

			// Count parentheses
			for _, ch := range trimmed {
				if ch == '(' {
					parenDepth++
				} else if ch == ')' {
					parenDepth--
					if parenDepth == 0 && inDefn {
						defn := currentDefn.String()
						if strings.HasPrefix(strings.TrimSpace(defn), "(defn") {
							defns = append(defns, defn)
						}
						currentDefn.Reset()
						inDefn = false
						break
					}
				}
			}
		}
	}

	return defns
}

// ExtractName extracts the function name from a defn form
func ExtractNameSimple(defn string) string {
	re := regexp.MustCompile(`\(defn\s+([^\s\[\(]+)`)
	matches := re.FindStringSubmatch(defn)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown"
}

// TestXsofyCountDefns just counts defns and tests bytecode compilation
func TestXsofyCountDefns(t *testing.T) {
	basePath := "/Users/ndn/development/xsofy/xsofy"

	// Collect all .lg files
	var lgFiles []string
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && strings.HasSuffix(path, ".lg") {
			lgFiles = append(lgFiles, path)
		}
		return nil
	})

	sort.Strings(lgFiles)

	totalDefns := 0
	bytecodeLowered := 0
	bytecodeErrors := 0

	for _, file := range lgFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Logf("read %s: %v", file, err)
			continue
		}

		defns := ExtractDefnsSimple(string(content))
		totalDefns += len(defns)

		for _, defn := range defns {
			// Try to compile via bytecode (default path)
			consts := vm.NewConsts()
			ns := rt.NS(rt.NameCoreNS)
			c := compiler.NewCompiler(consts, ns)
			c.SetSource("<xsofy-test>")

			_, _, err := c.CompileMultiple(strings.NewReader(defn))
			if err != nil {
				bytecodeErrors++
			} else {
				bytecodeLowered++
			}
		}
	}

	t.Logf("Xsofy defn analysis:\n")
	t.Logf("  Total defns: %d\n", totalDefns)
	t.Logf("  Bytecode compiled: %d (%.1f%%)\n", bytecodeLowered, 100.0*float64(bytecodeLowered)/float64(totalDefns))
	t.Logf("  Bytecode failed: %d (%.1f%%)\n", bytecodeErrors, 100.0*float64(bytecodeErrors)/float64(totalDefns))

	// This test is informational, don't fail
	if totalDefns == 0 {
		t.Error("Expected to find xsofy defns")
	}
}

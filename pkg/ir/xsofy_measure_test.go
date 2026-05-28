//go:build bootstrap

package ir_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type TestResult struct {
	Status  string // "lowered", "fallback", "error", "build-fail"
	Reason  string // fallback reason or error message
	File    string
	DefnNum int
	Name    string
}

// ExtractDefns extracts all top-level (defn ...) forms from Lisp source
func ExtractDefns(source string) []string {
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
func ExtractName(defn string) string {
	re := regexp.MustCompile(`\(defn\s+([^\s\[\(]+)`)
	matches := re.FindStringSubmatch(defn)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown"
}

func TestXsofyLowerGo(t *testing.T) {
	ensureLoader()

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
	t.Logf("Measuring lower-go pass rate against xsofy corpus\n")
	t.Logf("Found %d .lg files\n", len(lgFiles))

	var results []TestResult
	totalDefns := 0
	startTime := time.Now()
	tested := 0

	for _, file := range lgFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Logf("read %s: %v", file, err)
			continue
		}

		defns := ExtractDefns(string(content))
		relPath, _ := filepath.Rel(basePath, file)
		totalDefns += len(defns)

		for i, defn := range defns {
			name := ExtractName(defn)

			// Try to build and lower
			status, reason := "unknown", ""

			irFn := buildLispIR(t, defn)
			if irFn == nil {
				status = "build-fail"
				reason = "failed to build IR"
			} else {
				// Try to lower via ir.lower-go
				passVarCounter++
				varName := fmt.Sprintf("*test-ir-fn-%d*", passVarCounter)
				rt.NS(rt.NameCoreNS).Def(varName, irFn)

				expr := fmt.Sprintf(`(ir.lower-go/lower %s :bridge)`, varName)
				result := runLispExpr(t, expr)

				if result == nil {
					status = "error"
					reason = "nil result"
				} else if mapVal, ok := result.(*vm.PersistentMap); ok {
					statusVal := mapVal.ValueAt(vm.Keyword("status"))
					if statusVal != vm.NIL {
						status = fmt.Sprintf("%v", statusVal)
						if reasonVal := mapVal.ValueAt(vm.Keyword("reason")); reasonVal != vm.NIL {
							reason = fmt.Sprintf("%v", reasonVal)
						}
					} else {
						status = "unknown"
					}
				} else {
					status = "error"
					reason = fmt.Sprintf("unexpected result type: %T", result)
				}
			}

			results = append(results, TestResult{
				Status:  status,
				Reason:  reason,
				File:    relPath,
				DefnNum: i + 1,
				Name:    name,
			})
			tested++

			// Print progress every 50 defns
			if tested%50 == 0 {
				elapsed := time.Since(startTime).Seconds()
				rate := float64(tested) / elapsed
				remaining := int((float64(totalDefns) - float64(tested)) / rate)
				t.Logf("  Progress: %d/%d (%.1f defns/sec, ~%d sec remaining)\n", tested, totalDefns, rate, remaining)
			}
		}
	}

	elapsed := time.Since(startTime).Seconds()
	t.Logf("Completed in %.1f seconds\n", elapsed)

	// Aggregate statistics
	statusCounts := make(map[string]int)
	reasonCounts := make(map[string]int)

	for _, r := range results {
		statusCounts[r.Status]++
		if (r.Status == "fallback" || r.Status == "build-fail") && r.Reason != "" {
			reasonCounts[r.Reason]++
		}
	}

	t.Logf("Results:\n")
	t.Logf("  Total defns tested: %d\n", totalDefns)
	if totalDefns > 0 {
		t.Logf("  Lowered:    %d (%.1f%%)\n", statusCounts["lowered"], 100.0*float64(statusCounts["lowered"])/float64(totalDefns))
		t.Logf("  Fallback:   %d (%.1f%%)\n", statusCounts["fallback"], 100.0*float64(statusCounts["fallback"])/float64(totalDefns))
		t.Logf("  Build-Fail: %d (%.1f%%)\n", statusCounts["build-fail"], 100.0*float64(statusCounts["build-fail"])/float64(totalDefns))
		t.Logf("  Error:      %d (%.1f%%)\n", statusCounts["error"], 100.0*float64(statusCounts["error"])/float64(totalDefns))
	}

	if len(reasonCounts) > 0 {
		t.Logf("\nTop failure reasons:\n")

		// Sort reasons by frequency
		type reasonFreq struct {
			reason string
			count  int
		}
		var reasons []reasonFreq
		for reason, count := range reasonCounts {
			reasons = append(reasons, reasonFreq{reason, count})
		}
		sort.Slice(reasons, func(i, j int) bool {
			return reasons[i].count > reasons[j].count
		})

		for i, rf := range reasons {
			if i >= 10 {
				break
			}
			t.Logf("  %d. %s (%d)\n", i+1, rf.reason, rf.count)
		}
	}

	// Write results to file for inspection
	f, err := os.Create("/tmp/lower_go_results.txt")
	if err == nil {
		defer f.Close()
		w := bufio.NewWriter(f)
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.File, r.Name, r.Status, r.Reason)
		}
		w.Flush()
		t.Logf("\nDetailed results written to /tmp/lower_go_results.txt\n")
	}

	// Fail to show output
	t.Fail()
}

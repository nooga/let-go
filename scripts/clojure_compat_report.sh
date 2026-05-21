#!/usr/bin/env bash
set -euo pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
LOG=${TMPDIR:-/tmp}/let-go-clojure-compat-report.$$.$RANDOM.log

cd "$ROOT"

go test ./test/ -count=1 -run TestClojureTestSuite -v 2>&1 | tee "$LOG"
STATUS=$?

awk -v status="$STATUS" -v logfile="$LOG" '
function pct(n, d) {
	if (d == 0) {
		return "0.0"
	}
	return sprintf("%.1f", (n * 100.0) / d)
}

function add_known(name, pass, fail, error) {
	known_count++
	known_name[known_count] = name
	known_pass[known_count] = pass
	known_fail[known_count] = fail
	known_error[known_count] = error
}

function add_unexpected(name, kind, details) {
	unexpected_count++
	unexpected_name[unexpected_count] = name
	unexpected_kind[unexpected_count] = kind
	unexpected_details[unexpected_count] = details
}

function add_skip(name, reason) {
	skip_count++
	skip_name[skip_count] = name
	skip_reason[skip_count] = reason
}

function extract_counter(line, key,    pattern, value) {
	pattern = key "=[0-9]+"
	if (match(line, pattern)) {
		value = substr(line, RSTART + length(key) + 1, RLENGTH - length(key) - 1)
		return value + 0
	}
	return 0
}

/^=== RUN[[:space:]]+TestClojureTestSuite\// {
	current = $3
	sub(/^TestClojureTestSuite\//, "", current)
	next
}

/known failing/ {
	add_known(current, extract_counter($0, "pass"), extract_counter($0, "fail"), extract_counter($0, "error"))
	next
}

/FAILED/ {
	if (!(current in unexpected_seen)) {
		add_unexpected(current, "FAILED", $0)
		unexpected_seen[current] = 1
	}
	next
}

/PASSES but is listed in knownFailing/ {
	if (!(current in unexpected_seen)) {
		add_unexpected(current, "GRADUATED", $0)
		unexpected_seen[current] = 1
	}
	next
}

/--- FAIL: TestClojureTestSuite\// {
	name = $3
	sub(/^TestClojureTestSuite\//, "", name)
	if (!(name in unexpected_seen)) {
		add_unexpected(name, "FAIL", "go test marked this subtest failed")
		unexpected_seen[name] = 1
	}
	next
}

/--- SKIP: TestClojureTestSuite\// {
	name = $3
	sub(/^TestClojureTestSuite\//, "", name)
	skip_seen[name] = 1
	next
}

/compile:/ && current != "" {
	add_skip(current, "compile")
	next
}

/timeout after/ && current != "" {
	add_skip(current, "runtime timeout")
	next
}

/memory limit exceeded/ && current != "" {
	add_skip(current, "runtime memory")
	next
}

/panic:/ && current != "" {
	add_skip(current, "panic")
	next
}

/TOTALS:/ {
	files = extract_counter($0, "files")
	pass = extract_counter($0, "pass")
	fail = extract_counter($0, "fail")
	compile_skip = extract_counter($0, "compile")
	panic_skip = extract_counter($0, "panic")
	runtime_skip = extract_counter($0, "runtime")
	total = pass + fail
	next
}

END {
	print "Clojure compatibility suite report"
	print "==================================="
	if (files == 0) {
		print "No TOTALS line found. Full log: " logfile
		exit status
	}

	printf("Assertions: %d / %d passing (%s%%)\n", pass, total, pct(pass, total))
	printf("Failures:   %d assertions\n", fail)
	printf("Files:      %d run, %d skipped (compile=%d panic=%d runtime=%d)\n", files, compile_skip + panic_skip + runtime_skip, compile_skip, panic_skip, runtime_skip)
	printf("Go status:  %d\n", status)
	print ""

	if (unexpected_count > 0) {
		printf("Unexpected failing files (%d):\n", unexpected_count)
		for (i = 1; i <= unexpected_count; i++) {
			printf("  - %s [%s]\n", unexpected_name[i], unexpected_kind[i])
		}
		print ""
	} else {
		print "Unexpected failing files: none"
		print ""
	}

	if (skip_count > 0) {
		printf("Skipped files (%d):\n", skip_count)
		for (i = 1; i <= skip_count; i++) {
			printf("  - %s (%s)\n", skip_name[i], skip_reason[i])
		}
		print ""
	}

	if (known_count > 0) {
		printf("Known failing files (%d):\n", known_count)
		for (i = 1; i <= known_count; i++) {
			printf("  - %s: pass=%d fail=%d error=%d\n", known_name[i], known_pass[i], known_fail[i], known_error[i])
		}
		print ""
	}

	print "Full verbose log: " logfile
	exit status
}
' "$LOG"

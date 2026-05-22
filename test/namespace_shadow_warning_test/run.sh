#!/bin/bash
set -e
cd "$(dirname "$0")/../.."

LG="./lg"
DIR="test/namespace_shadow_warning_test"
fail=0

check() {
  local file="$1"
  local should_have="$2"
  local err
  err=$($LG "$file" 2>&1 >/dev/null)
  if [ "$should_have" = "warn" ]; then
    if echo "$err" | grep -q "already refers to: #'clojure.core/"; then
      echo "  PASS: $file emitted warning"
    else
      echo "  FAIL: $file should have emitted warning. stderr: '$err'"
      fail=$((fail+1))
    fi
  else
    if echo "$err" | grep -q "already refers to: #'clojure.core/"; then
      echo "  FAIL: $file should NOT have emitted warning. stderr: '$err'"
      fail=$((fail+1))
    else
      echo "  PASS: $file silent"
    fi
  fi
}

check "$DIR/01_shadow_emits_warning.lg"      warn
check "$DIR/02_exclude_suppresses.lg"        no-warn
check "$DIR/03_non_core_name_no_warning.lg"  no-warn
check "$DIR/04_local_binding_no_warning.lg"  no-warn
check "$DIR/05_redef_own_var_no_warning.lg"  no-warn

exit $fail

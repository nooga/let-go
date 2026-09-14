#!/usr/bin/env bash
# ci-usage.sh — where this repo's GitHub Actions runner time actually goes.
#
# The Actions usage-metrics REST API is org/enterprise-scoped, and this repo
# is owned by a user account, so `repos/:owner/:repo/actions/metrics/usage`
# 404s and the /actions/metrics/usage page has no programmatic equivalent.
# Billing endpoints need the owner's own token. What is left, and what this
# uses, is the ordinary runs + jobs API.
#
# That turns out to be strictly better for diagnosis, because it exposes
# three things the metrics page does not:
#
#   1. Queue wait (job.created_at -> job.started_at). This is what separates
#      "we are over an account-wide concurrency cap" from "one workflow is
#      serializing itself behind its own concurrency: group". The two look
#      identical on the metrics page and have opposite fixes. If waits are
#      spread across every workflow it is the former; if they are confined to
#      one workflow it is that workflow's own group.
#   2. Skipped jobs, which inflate the run count on the metrics page while
#      costing nothing. Several perf workflows here skip every job on most
#      PRs, so raw run counts badly overstate load.
#   3. Jobs cancelled at exactly their timeout-minutes ceiling — compute
#      spent producing no artifact. Reported separately under "wasted",
#      because it is the one line item that is pure loss.
#
# Usage: scripts/ci-usage.sh [owner/repo] [since-YYYY-MM-DD]
#
# Defaults to the repo of the current checkout and the last 7 days. Needs
# `gh` (authenticated) and `jq`. Roughly one API call per run: a 7-day
# window here is ~700 calls against a 5000/hr limit.
set -euo pipefail

REPO="${1:-$(gh repo view --json nameWithOwner --jq .nameWithOwner)}"
SINCE="${2:-$(date -v-7d +%Y-%m-%d 2>/dev/null || date -d '7 days ago' +%Y-%m-%d)}"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

echo "repo=$REPO since=$SINCE" >&2

gh api --paginate "repos/$REPO/actions/runs?created=>=$SINCE&per_page=100" \
  --jq '.workflow_runs[] | {id,name,event,conclusion,created_at,head_branch}' > "$WORK/runs.jsonl"
RUNS=$(wc -l < "$WORK/runs.jsonl" | tr -d ' ')

# The runs list caps at 1000 results however you paginate. Say so rather than
# silently reporting a truncated window as if it were the whole thing.
if [ "$RUNS" -ge 1000 ]; then
  echo "WARNING: hit the 1000-run pagination cap — narrow the window to see further back" >&2
fi

# filter=latest keeps re-run attempts from being counted twice.
cat > "$WORK/one.sh" <<EOF
#!/bin/bash
gh api "repos/$REPO/actions/runs/\$1/jobs?per_page=100&filter=latest" \
  --jq '.jobs[] | {wf:.workflow_name, job:.name, conclusion, created_at, started_at, completed_at, labels:(.labels|join(","))}' 2>/dev/null
EOF
chmod +x "$WORK/one.sh"
jq -r '.id' "$WORK/runs.jsonl" | xargs -P 8 -n 1 "$WORK/one.sh" > "$WORK/jobs.jsonl"

jq -c 'select(.started_at != null and .completed_at != null)
  | .dur   = ((.completed_at|fromdateiso8601) - (.started_at|fromdateiso8601))
  | .queue = ((.started_at|fromdateiso8601)   - (.created_at|fromdateiso8601))
  | .os = (if   (.labels|test("macos"))   then "macos"
           elif (.labels|test("windows")) then "windows"
           else "ubuntu" end)' "$WORK/jobs.jsonl" > "$WORK/d.jsonl"

JOBS=$(wc -l < "$WORK/d.jsonl" | tr -d ' ')
echo
echo "=== $REPO since $SINCE: $RUNS runs, $JOBS jobs ==="

echo
echo "=== runner-minutes by workflow (skipped jobs excluded) ==="
jq -s '[.[]|select(.conclusion!="skipped")]|group_by(.wf)
  |map({wf:.[0].wf, jobs:length, min:((map(.dur)|add)/60|round)})|sort_by(-.min)[]
  |"\(.min)\tmin\t\(.jobs)\tjobs\t\(.wf)"' -r "$WORK/d.jsonl"

echo
echo "=== runner-minutes by runner OS ==="
jq -s '[.[]|select(.conclusion!="skipped")]|group_by(.os)
  |map({os:.[0].os, jobs:length, min:((map(.dur)|add)/60|round)})|sort_by(-.min)[]
  |"\(.min)\tmin\t\(.jobs)\tjobs\t\(.os)"' -r "$WORK/d.jsonl"

# Matrix legs are folded together: "profile (macos-14)" and "profile
# (ubuntu-latest)" report as one "profile" row, since the question here is
# which step costs, not which leg.
echo
echo "=== top jobs by total minutes (matrix legs folded) ==="
jq -s '[.[]|select(.conclusion!="skipped")]
  |group_by(.wf+" :: "+(.job|gsub("\\(.*";"")))
  |map({k:(.[0].wf+" :: "+(.[0].job|gsub("\\(.*";""))), n:length,
        min:((map(.dur)|add)/60|round), avg:((map(.dur)|add)/length/60*10|round/10)})
  |sort_by(-.min)[:15][]|"\(.min)\tmin\tn=\(.n)\tavg=\(.avg)m\t\(.k)"' -r "$WORK/d.jsonl"

echo
echo "=== wasted: minutes burned by jobs that ended cancelled or failed ==="
jq -s '[.[]|select(.conclusion=="cancelled" or .conclusion=="failure")]|group_by(.wf+"/"+.os)
  |map({k:(.[0].wf+" / "+.[0].os), n:length, min:((map(.dur)|add)/60|round)})|sort_by(-.min)[]
  |"\(.min)\tmin\tn=\(.n)\t\(.k)"' -r "$WORK/d.jsonl"

echo
echo "=== queue wait: created -> started ==="
echo "=== waits across every workflow => account concurrency cap"
echo "=== waits confined to one workflow => that workflow's own concurrency: group"
jq -s '[.[]|select(.conclusion!="skipped")]|group_by(.wf+"/"+.os)
  |map({k:(.[0].wf+" / "+.[0].os), n:length,
        median:(((map(.queue)|sort)[(length/2|floor)])/60|round),
        max:((map(.queue)|max)/60|round),
        over10:([.[]|select(.queue>600)]|length)})|sort_by(-.max)[]
  |"\(.k)\tn=\(.n)\tmedian=\(.median)m\tmax=\(.max)m\t>10m:\(.over10)"' -r "$WORK/d.jsonl"

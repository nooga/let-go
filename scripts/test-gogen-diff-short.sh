#!/bin/sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp_dir=$(mktemp -d "${TMPDIR:-/tmp}/letgo-gogen-diff-short.XXXXXX")
trap 'rm -rf "$tmp_dir"' EXIT HUP INT TERM

probe_file="$tmp_dir/probe"
fake_go="$tmp_dir/go"
cat >"$fake_go" <<'FAKE_GO'
#!/bin/sh
set -eu

record_leg() {
	leg=$1
	case " $* " in
		*" -short=false "*) printf '%s\n' "$leg" >>"$PARITY_PROBE" ;;
		*) printf '%s:skipped\n' "$leg" >>"$PARITY_PROBE" ;;
	esac
}

case "$*" in
	*TestParityGatePhase1*) record_leg parity "$@" ;;
	*TestGogenAOTDiff*) record_leg gogen-diff "$@" ;;
esac
exit 0
FAKE_GO
chmod +x "$fake_go"

PATH="$tmp_dir:$PATH" GOFLAGS=-short PARITY_PROBE="$probe_file" \
	make -s -C "$repo_root" gogen-diff >/dev/null

expected="$tmp_dir/expected"
printf 'parity\ngogen-diff\n' >"$expected"
if ! cmp -s "$expected" "$probe_file"; then
	echo "gogen-diff must force -short=false for both mandatory test legs" >&2
	echo "observed:" >&2
	sed 's/^/  /' "$probe_file" >&2
	exit 1
fi

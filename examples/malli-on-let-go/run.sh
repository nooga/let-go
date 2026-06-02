#!/usr/bin/env bash
#
# Launch an interactive let-go REPL with metosin/malli loaded.
#
# lgx.edn is the dependency manifest. lgx's own runner can't attach an
# interactive REPL yet (it shells out via buffered os/sh with a dead stdin —
# see its TODO; the fix needs os/exec*, added in a separate PR). So this interim
# launcher reads lgx.edn to build -source-paths and hands them to `lg -r`. Once
# lgx's runner uses os/exec*, this whole script collapses to `lgx run`.
#
# Quit the REPL with Ctrl-C (or Ctrl-D).
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$HERE/../.." && pwd)"          # the let-go repo (carries the malli fixes)
LG="${LG:-$HERE/lg-example}"               # patched let-go binary; override with LG=...

# Build the let-go binary once. The example needs the compiler/VM fixes in this
# repo (PRs #144-#153) — stock `lg` won't run malli. `rm "$LG"` forces a rebuild.
if [ ! -x "$LG" ]; then
  echo "building let-go (one-time) ..." >&2
  ( cd "$REPO" && go build -o "$LG" . )
fi

# Resolve -source-paths from lgx.edn: this project's :paths plus each dep's
# :local/root (+ optional :deps/root). Uses `lg` itself to read the edn.
PATHS="$("$LG" -e "(do (print (let [cfg  (read-string (slurp \"$HERE/lgx.edn\"))
                                    own  (map (fn [p] (str \"$HERE/\" p)) (:paths cfg))
                                    deps (map (fn [[_ c]] (str (:local/root c)
                                                              (if (:deps/root c) (str \"/\" (:deps/root c)) \"\")))
                                              (:deps cfg))]
                                (apply str (interpose \":\" (concat own deps))))) (os/exit 0))")"

# LG_READ_BB / LG_READ_CLJ enable malli's #?(:bb ...) / #?(:clj ...) reader
# conditionals. -r attaches the REPL after :main (repl.lg) runs.
echo "starting malli REPL (Ctrl-C to quit) ..." >&2
exec env LG_READ_BB=1 LG_READ_CLJ=1 "$LG" -r -source-paths "$PATHS" "$HERE/repl.lg"

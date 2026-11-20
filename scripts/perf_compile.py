#!/usr/bin/env python3
"""A/B readout for programs built by `lg compile` — informational.

The perf-pr checks time Go micro-benchmarks, which cannot see what lowering a
program buys or costs. This builds each program under benchmark/compile three
ways and times the results:

  vm        lg <file>                      the bytecode VM
  compiled  lg compile                     the program lowered to Go
  core      lg compile -lowered-core       clojure.core and clojure.walk too

Each program times its own workload and prints one line,
`PERF <name> <best microseconds> <checksum>`, so process start and the build
are not in the number.

  perf_compile.py build --lg BIN --src CHECKOUT --programs DIR --out DIR
  perf_compile.py run --base DIR --head DIR [--rounds N]

`build` is run once per side, each with that side's lg and checkout but the
SAME programs directory, so a PR that edits a program still compares like with
like. `run` interleaves the two sides and prints a markdown report.

It never fails on a timing: shared runners are too noisy for a few percent to
mean anything. Read the count of lowered functions first (a drop means code
went back to the VM) and the ratios between arms second. The exit code is
nonzero only when the measurement itself is broken: a head build or run that
fails, or an arm whose checksum differs between base and head. Arms that
disagree with each other at head are listed as a note, since a compiled
program computing something else than the VM is worth seeing but is not this
PR's doing unless the count changed.
"""
import argparse
import json
import os
import re
import subprocess
import sys

ARMS = ("vm", "compiled", "core")
LOWERED_RE = re.compile(r"(\d+) fns lowered")


def programs_in(directory):
    return sorted(f[:-3] for f in os.listdir(directory) if f.endswith(".lg"))


def build(args):
    os.makedirs(args.out, exist_ok=True)
    env = dict(os.environ, LETGO_SRC=os.path.abspath(args.src), LG_SOURCE_PATHS=os.path.abspath(args.programs))
    lg = os.path.abspath(args.lg)
    meta = {}
    for name in programs_in(args.programs):
        source = os.path.join(os.path.abspath(args.programs), name + ".lg")
        entry = {"vm": {"cmd": [lg, source]}}
        for arm, flags in (("compiled", []), ("core", ["-lowered-core"])):
            binary = os.path.join(os.path.abspath(args.out), f"{name}.{arm}")
            proc = subprocess.run([lg, "compile", *flags, "-o", binary, source], env=env, capture_output=True, text=True)
            if proc.returncode != 0:
                # A base that predates a flag has no such arm; the report says so.
                entry[arm] = {"error": (proc.stdout + proc.stderr).strip().splitlines()[-1:]}
                continue
            lowered = LOWERED_RE.search(proc.stdout)
            entry[arm] = {"cmd": [binary], "lowered": int(lowered.group(1)) if lowered else None}
        meta[name] = entry
        print(f"built {name}: " + ", ".join(f"{a}={'ok' if 'cmd' in entry[a] else 'FAILED'}" for a in ARMS), file=sys.stderr)
    with open(os.path.join(args.out, "meta.json"), "w") as f:
        json.dump(meta, f, indent=1)


def measure(cmd, name):
    """One run: (microseconds, checksum), or (None, reason)."""
    proc = subprocess.run(cmd, capture_output=True, text=True)
    for line in proc.stdout.splitlines():
        parts = line.split()
        if len(parts) == 4 and parts[0] == "PERF" and parts[1] == name:
            return int(parts[2]), parts[3]
    return None, (proc.stderr.strip().splitlines() or [f"exit {proc.returncode}, no PERF line"])[-1]


def load(directory):
    with open(os.path.join(directory, "meta.json")) as f:
        return json.load(f)


def fmt_us(us):
    return "n/a" if us is None else f"{us / 1000:.1f} ms"


def run(args):
    sides = {"base": load(args.base), "head": load(args.head)}
    names = sorted(sides["head"])
    best, sums, broken, notes = {}, {}, [], []
    for rnd in range(args.rounds):
        # ABBA: alternate which side goes first, so slow drift on the runner
        # does not land on one side.
        order = ("base", "head") if rnd % 2 == 0 else ("head", "base")
        for name in names:
            for arm in ARMS:
                for side in order:
                    cmd = sides[side].get(name, {}).get(arm, {}).get("cmd")
                    if not cmd:
                        continue
                    us, extra = measure(cmd, name)
                    key = (side, name, arm)
                    if us is None:
                        if side == "head":
                            broken.append(f"{name} {arm} at head did not run: {extra}")
                        continue
                    best[key] = min(us, best.get(key, us))
                    sums.setdefault(key, set()).add(extra)
    for name in names:
        for arm in ARMS:
            if "cmd" not in sides["head"][name].get(arm, {}):
                broken.append(f"{name} {arm} at head did not build: {sides['head'][name].get(arm, {}).get('error')}")
            b, h = sums.get(("base", name, arm), set()), sums.get(("head", name, arm), set())
            if len(h) > 1 or (b and h and b != h):
                broken.append(f"{name} {arm}: checksum differs between base and head, or between runs: {sorted(b)} vs {sorted(h)}")
        at_head = {arm: sorted(sums.get(("head", name, arm), ())) for arm in ARMS}
        if len({tuple(v) for v in at_head.values() if v}) > 1:
            notes.append(f"`{name}`: arms disagree on the result at head: " + ", ".join(f"{a} {v[0]}" for a, v in at_head.items() if v))

    print(f"Best of {args.rounds} runs per arm, base and head interleaved on one runner. Informational, not a gate.\n")
    print("| Program | Arm | Base | Head | Delta | Program fns lowered (base → head) |")
    print("|---|---|---:|---:|---:|---|")
    for name in names:
        for arm in ARMS:
            b, h = best.get(("base", name, arm)), best.get(("head", name, arm))
            delta = "n/a" if not (b and h) else f"{(h / b - 1) * 100:+.1f}%"
            lb = sides["base"].get(name, {}).get(arm, {}).get("lowered")
            lh = sides["head"][name].get(arm, {}).get("lowered")
            lowered = "" if arm == "vm" else f"{'n/a' if lb is None else lb} → {'n/a' if lh is None else lh}"
            print(f"| `{name}` | {arm} | {fmt_us(b)} | {fmt_us(h)} | {delta} | {lowered} |")
    print("\nWhat lowering buys at head (speedup over the VM):\n")
    print("| Program | compiled | core |")
    print("|---|---:|---:|")
    for name in names:
        vm = best.get(("head", name, "vm"))
        cells = []
        for arm in ("compiled", "core"):
            h = best.get(("head", name, arm))
            cells.append("n/a" if not (vm and h) else f"{vm / h:.2f}×")
        print(f"| `{name}` | {cells[0]} | {cells[1]} |")
    if notes:
        print("\nNotes:\n")
        for line in notes:
            print(f"- {line}")
    if broken:
        print("\n**Measurement problems:**\n")
        for line in broken:
            print(f"- {line}")
    return 1 if broken else 0


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = parser.add_subparsers(dest="command", required=True)
    b = sub.add_parser("build")
    b.add_argument("--lg", required=True)
    b.add_argument("--src", required=True)
    b.add_argument("--programs", required=True)
    b.add_argument("--out", required=True)
    r = sub.add_parser("run")
    r.add_argument("--base", required=True)
    r.add_argument("--head", required=True)
    r.add_argument("--rounds", type=int, default=5)
    args = parser.parse_args()
    if args.command == "build":
        build(args)
        return 0
    return run(args)


if __name__ == "__main__":
    sys.exit(main())

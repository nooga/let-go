#!/usr/bin/env python3
"""Interleaved repeated A/B for perf-pr variance reduction — informational.

Runs N interleaved base/head benchmark snapshots across two pre-built worktrees
and reports the per-family delta distribution. It NEVER fails the job: the exit
code reflects harness integrity (see below), not a perf verdict. The point of
this phase is to collect the data needed to calibrate a gate, not to be one.

Why repeat + interleave: a single-shot base-vs-head A/B on shared CI runners is
too heavy-tailed to gate on — ~1 in 4 runs blows a memory-bound cluster to
25-40% while the register-only anchor stays flat (nooga/let-go#445). Repetition
lets the per-family MEDIAN absorb up to floor((N-1)/2) contaminated cycles;
interleaving cancels slow time-separation drift.

Each side is snapshotted with `bench-ratchet -profile <p> ... snapshot`, which
emits ratio_to_anchor per benchmark; delta = head_ratio / base_ratio - 1.

Reported per run (a same-code no-op PR should show zero would-gates):
  - median-of-N delta per family, and whether ANY family would gate at each
    candidate budget — the decision unit is the whole workflow, so we report
    "would any family gate", one Bernoulli sample toward P(any-gate | no-op).
  - confirm K-of-N as a secondary statistic (dominated within one job, since
    contaminated cycles cluster — kept only for corroboration).

Integrity (drives the exit code, so a broken measurement can't read as clean):
  - MISSING / NEW families: present on every cycle of one side but never the
    other. A rename/add/remove between base and head — reported, not silently
    dropped (the old `set(base) & set(head)` hid these).
  - FLAKY families: present on some cycles of a side but not all, i.e. not
    exactly N observations. That means the two halves aren't comparable, which
    invalidates the deltas — exit nonzero.
"""
import argparse, json, os, statistics, subprocess, sys


def snapshot(worktree, profile, out, timeout, label):
    """Run one bench-ratchet snapshot in a worktree.

    label is a human tag ("base cycle 3/7") so a failed `go run` or an
    unexpected bench-ratchet JSON shape surfaces as one triage line instead
    of a raw traceback. The failure still exits nonzero — a missing sample is
    a harness-integrity failure, not a perf verdict — just legibly.

    Returns (ratios, samples, meta):
      ratios  {family: ratio_to_anchor}
      samples {family: [ratio_to_anchor, ...]}   raw per-sample distribution
      meta    {anchor_ns, cpu_model, go_version, num_cpu, arch, sha}
    """
    try:
        subprocess.run(
            ["go", "run", "./cmd/bench-ratchet", "-profile", profile,
             "-timeout", timeout, "-baseline", out, "snapshot"],
            cwd=worktree, check=True,
        )
        with open(out) as f:
            d = json.load(f)
        (_, m), = d["machines"].items()
    except (subprocess.CalledProcessError, OSError,
            json.JSONDecodeError, KeyError, ValueError) as e:
        print(f"::error::{label} failed: {type(e).__name__}: {e}", flush=True)
        raise
    ratios = {k: e["ratio_to_anchor"] for k, e in m["benchmarks"].items()}
    samples = {k: [s.get("ratio_to_anchor") for s in e.get("samples", [])]
               for k, e in m["benchmarks"].items()}
    mach = m.get("machine", {})
    meta = {
        "anchor_ns": m["anchor"]["ns_per_op"],
        "cpu_model": mach.get("cpu_model", "?"),
        "go_version": mach.get("go_version", "?"),
        "num_cpu": mach.get("num_cpu", 0),
        "arch": mach.get("arch", "?"),
        "sha": m.get("captured_at_sha", "?"),
    }
    return ratios, samples, meta


def short(fam):
    return fam.replace("github.com/nooga/let-go/pkg/vm.Benchmark", "")


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--base", required=True, help="pre-built base worktree")
    ap.add_argument("--head", required=True, help="pre-built head worktree")
    ap.add_argument("--n", type=int, default=7, help="interleaved cycle count")
    ap.add_argument("--profile", default="pr-fast")
    ap.add_argument("--budgets", default="6,8,10",
                    help="comma-separated candidate budgets %% to report would-gate for")
    ap.add_argument("--confirm-k", type=int, default=2,
                    help="cycles that must exceed the primary budget (secondary stat)")
    ap.add_argument("--timeout", default="15m")
    ap.add_argument("--out", default="ab-out")
    args = ap.parse_args()
    os.makedirs(args.out, exist_ok=True)
    budgets = [float(b) for b in args.budgets.split(",")]
    primary = budgets[-1]  # confirm-K reported at the widest (most lenient) budget

    # Per-side, per-family: list of ratios, one entry per cycle it appeared in.
    base_seen, head_seen = {}, {}
    base_samples, head_samples = {}, {}
    anchors = []          # (base_ns, head_ns) per cycle
    order = []            # measurement order per cycle: "BH" or "HB"
    machines = {"base": None, "head": None}

    for i in range(1, args.n + 1):
        print(f"::group::cycle {i}/{args.n}", flush=True)
        # Counterbalance (ABBA): odd cycles base-first, even head-first. Base and
        # head are identical code on a no-op PR, so a fixed base-then-head order
        # lets first-vs-second-position drift (warmup, cache, thermal) look like a
        # consistent head regression — a systematic bias the median cannot remove.
        # Alternating cancels it. (Surfaced porting to paserati: a fixed order
        # produced a 1-in-N false positive on a no-op PR.)
        base_first = (i % 2 == 1)
        order.append("BH" if base_first else "HB")

        def do_base():
            r, s, meta = snapshot(args.base, args.profile,
                                  f"{args.out}/base_{i}.json", args.timeout,
                                  f"base cycle {i}/{args.n}")
            machines["base"] = machines["base"] or meta
            for fam, v in r.items():
                base_seen.setdefault(fam, []).append(v)
                base_samples.setdefault(fam, []).append(s.get(fam))
            return meta["anchor_ns"]

        def do_head():
            r, s, meta = snapshot(args.head, args.profile,
                                  f"{args.out}/head_{i}.json", args.timeout,
                                  f"head cycle {i}/{args.n}")
            machines["head"] = machines["head"] or meta
            for fam, v in r.items():
                head_seen.setdefault(fam, []).append(v)
                head_samples.setdefault(fam, []).append(s.get(fam))
            return meta["anchor_ns"]

        if base_first:
            ba = do_base()
            ha = do_head()
        else:
            ha = do_head()
            ba = do_base()
        anchors.append((ba, ha))
        print("::endgroup::", flush=True)

    # --- Integrity: categorize families before computing any delta. ------------
    # A family is comparable only if it appears on EXACTLY N base and N head
    # cycles. Anything else is surfaced, not silently intersected away.
    base_fams, head_fams = set(base_seen), set(head_seen)
    comparable, missing, new, flaky = [], [], [], []
    for fam in sorted(base_fams | head_fams):
        nb, nh = len(base_seen.get(fam, [])), len(head_seen.get(fam, []))
        if nb == args.n and nh == args.n:
            comparable.append(fam)
        elif nb == args.n and nh == 0:
            missing.append(fam)      # dropped/renamed in head
        elif nb == 0 and nh == args.n:
            new.append(fam)          # added/renamed in head
        else:
            flaky.append((fam, nb, nh))  # inconsistent presence — harness noise

    rows = []
    for fam in comparable:
        ds = [(head_seen[fam][k] / base_seen[fam][k] - 1.0) * 100
              for k in range(args.n) if base_seen[fam][k]]
        med = statistics.median(ds)
        rows.append({
            "fam": short(fam),
            "n": len(ds),
            "median": med,
            "worst": max(ds, key=abs),
            "deltas": ds,
            "exceed_primary": sum(1 for d in ds if d > primary),
            "base_samples": base_samples.get(fam),
            "head_samples": head_samples.get(fam),
        })
    rows.sort(key=lambda r: r["median"], reverse=True)

    # would-gate at each budget: does ANY comparable family's median exceed it?
    would_gate = {}
    for bud in budgets:
        hits = [r["fam"] for r in rows if r["median"] > bud]
        would_gate[bud] = hits
    confirm_hits = [r["fam"] for r in rows
                    if r["exceed_primary"] >= args.confirm_k]

    aggregate = {
        "provenance": {
            "n": args.n, "profile": args.profile, "budgets": budgets,
            "confirm_k": args.confirm_k,
            "samples_per_snapshot": len(next(iter(base_samples.values()), []) or []),
            "cycle_order": order,
            "anchors": anchors,
            "machine_base": machines["base"],
            "machine_head": machines["head"],
        },
        "integrity": {
            "comparable": len(comparable),
            "missing": [short(f) for f in missing],
            "new": [short(f) for f in new],
            "flaky": [[short(f), nb, nh] for f, nb, nh in flaky],
        },
        "would_gate": {str(b): hits for b, hits in would_gate.items()},
        "confirm_hits": confirm_hits,
        "rows": rows,
    }
    with open(f"{args.out}/aggregate.json", "w") as f:
        json.dump(aggregate, f, indent=2)

    # --- Human report (job summary). ------------------------------------------
    mb = machines["base"] or {}
    print(f"\n=== interleaved repeat A/B — N={args.n}, profile={args.profile} "
          f"(INFORMATIONAL, never gates) ===")
    print(f"runner: {mb.get('cpu_model','?')}  "
          f"{mb.get('num_cpu','?')} CPU  {mb.get('go_version','?')}")
    print(f"base {(machines['base'] or {}).get('sha','?')[:12]}  "
          f"head {(machines['head'] or {}).get('sha','?')[:12]}  "
          f"order {' '.join(order)}")
    print("anchor stability per cycle (base/head ns/op, "
          "should stay flat — it measures compute, not memory):")
    for i, (ba, ha) in enumerate(anchors, 1):
        drift = (ha / ba - 1) * 100 if ba else float("nan")
        print(f"  cycle {i}: base {ba:.3f}  head {ha:.3f}  Δ{drift:+.1f}%")

    if missing or new:
        print(f"\nfamily-set change vs base:  "
              f"MISSING={[short(f) for f in missing] or '-'}  "
              f"NEW={[short(f) for f in new] or '-'}")
    if flaky:
        print(f"\n!! FLAKY families (not exactly N obs — measurement suspect): "
              f"{[(short(f), nb, nh) for f, nb, nh in flaky]}")

    print(f"\n{'family':40} {'median%':>8} {'worst%':>7} {'>prim':>5}")
    print("-" * 64)
    for r in rows[:12]:
        print(f"{r['fam']:40} {r['median']:+8.2f} {r['worst']:+7.2f} "
              f"{r['exceed_primary']:5d}")
    if len(rows) > 12:
        print(f"... {len(rows) - 12} more (full table in aggregate.json)")

    print("\nwould-gate (any comparable family's median > budget):")
    for bud in budgets:
        hits = would_gate[bud]
        print(f"  budget {bud:4.0f}%:  {'GATE' if hits else 'clean':5}  "
              f"{hits or ''}")
    print(f"confirm K={args.confirm_k}/{args.n} @ {primary:.0f}% (secondary): "
          f"{confirm_hits or 'clean'}")

    # On a no-op PR any would-gate is a FALSE POSITIVE; each budget's boolean is
    # one Bernoulli sample toward P(any-gate | no-op), aggregated across runs.
    print("\n(no-op PR: every would-gate above is a false positive; "
          "each is one sample toward P(any-gate).)")

    # Exit code = integrity only, never a perf verdict. Flaky presence means the
    # halves aren't comparable, so the run's deltas are not trustworthy.
    if flaky:
        print("\nFAIL: flaky family presence — deltas not trustworthy.",
              file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()

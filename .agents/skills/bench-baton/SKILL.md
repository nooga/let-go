---
name: bench-baton
description: Coordinate heavy local workloads across worktrees, processes, and subagents — benchmarks and timing-sensitive gates run exclusively on a quiesced machine, while builds, test suites, regeneration, and lint run concurrently in a capacity-limited shared lane. Use when parallel workers could perturb benchmark numbers, when a perf gate needs an idle machine, or when concurrent builds need a fair capacity cap and a shared ledger.
version: 2026-09-03
triggers: ["benchmark", "bench-ratchet", "perf gate", "quiet machine", "run tests in parallel", "make generate", "baton", "lease"]
tools: [bash]
preconditions: ["Unix-like host (flock)", "go toolchain"]
constraints: ["benchmarks and timing gates go through bench, never build", "give every worker a stable --owner", "cheap commands take no lease"]
---

# bench-baton — one quiet machine, two lanes

The scarce resource is a **quiet machine**, not the build tree. `cmd/bench-baton`
arbitrates one pool per machine-level resource:

| lane | mode | use for | concurrency |
|---|---|---|---|
| `bench` | `exclusive` | benchmarks, `make bench-ratchet`, timing assertions | alone — no other bench, no builds |
| `build` | `shared` | builds, test suites, `make generate` / `check-generated`, lint | up to `--max-shared` (default 4) |

A waiting benchmark **closes the gate**: shared work arriving behind it is deferred
until the bench has run, so a stream of builds cannot starve it. Leases are `flock`s,
so a holder that dies releases the pool by itself; `reap` only tidies the record.

State lives in `$BENCH_BATON_HOME` (default `~/Library/Caches/let-go/batons` on macOS,
`~/.cache/let-go/batons` on Linux), one directory per `--pool` (default `letgo`):
`state.json`, `ledger.jsonl`, `logs/`, and the lock files. Set `BENCH_BATON_HOME`
to an existing pool root to keep its ledger.

## Shell

```bash
go run ./cmd/bench-baton bench --owner pr641 --cwd WT -- make bench-ratchet
go run ./cmd/bench-baton build --owner pr641 --cwd WT -- go test ./...
go run ./cmd/bench-baton run --mode shared --owner pr641 -- make generate
go run ./cmd/bench-baton status
go run ./cmd/bench-baton ledger -n 10
go run ./cmd/bench-baton reap
```

`build` / `bench` exit with the wrapped command's return code and print its tail, so
they drop into existing scripts unchanged; the full output is in the log the ledger
names. `--lease-timeout SEC` makes a worker report contention instead of blocking;
`--timeout SEC` kills a runaway command.

## Rules

- Benchmarks and any timing-sensitive gate go through `bench` — never `build`. A
  ratchet number captured next to a running test suite is noise the ratchet cannot
  tell from a regression (the pre-push gate compares against a baseline captured on
  a quiet machine).
- Builds and tests use `build`; they are concurrent, so do not hand-serialize them.
- Give every worker a stable `--owner` (worktree or PR name) so `status` and `ledger`
  read well across subagents.
- One pool per machine-level resource; `--pool` for an unrelated project.
- Cheap commands (grep, one focused test, `jj status`) take no lease at all.
- Before blaming a commit for an anchor drift, check the power state
  (`pmset -g | grep lowpowermode`); a laptop on battery is not a quiet machine either.

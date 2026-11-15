# What `lg compile` buys

Four small programs, each built three ways and timed, so there is a standing
answer to "what does lowering this kind of code to Go buy":

| Program | Shape |
|---|---|
| `maps.lg` | core calling core over nested maps (`assoc-in`, `get-in`, `select-keys`) |
| `callbacks.lg` | core calling a user function per element (`every?`, `not-any?`, a `filter` transducer, `remove`) |
| `numeric.lg` | a hinted float loop with no collections |
| `page.lg` | text built from records; mostly natives and allocation |

| Arm | Build |
|---|---|
| `vm` | `lg <file>` |
| `compiled` | `lg compile` |
| `core` | `lg compile -lowered-core` |

Each program times its own workload, best of seven, and prints
`PERF <name> <microseconds> <checksum>`.

## Run it

```sh
go build -buildvcs=false -o /tmp/lg .
python3 scripts/perf_compile.py build --lg /tmp/lg --src . --programs benchmark/compile --out /tmp/perf-compile/head
python3 scripts/perf_compile.py run --base /tmp/perf-compile/head --head /tmp/perf-compile/head
```

That compares a build against itself, which shows the ratios between arms and
how much two identical runs differ on your machine. To compare two commits,
build each side with its own `lg` and checkout and pass both directories.

## In CI

The `perf-compile` label on a pull request runs the same comparison between
the PR's merge-base and its head and posts the report as a comment
(`.github/workflows/perf-compile.yml`). It is informational and never fails a
PR on a timing.

## Reading the report

- **Fns lowered** first. A count that drops means code went back to the VM,
  which a noisy timing can hide.
- **The ratios between arms** second. They move less than the absolute times.
- A few percent between base and head on a shared runner is noise.
- `numeric` can print a different checksum on the VM and compiled arms on
  arm64, where Go fuses `x*y+z` into one rounding. The report lists it as a
  note.

Adding a program: one file here with an `ns` form, a `-main` that prints the
`PERF` line, and the `(when-not *compiling-aot* (-main))` guard.

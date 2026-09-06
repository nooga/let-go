# AGENTS.md — let-go

let-go is a Clojure dialect in Go: reader → bytecode compiler → stack VM, no JVM, plus an
opt-in IR path and a Go AOT backend. The README covers the language and the embedding API.
Contributor material lives under `docs/`: start at `docs/README.md`, then
`docs/contributor-workflow.md` (commands, generated artifacts, CI gates) and
`docs/contribution-policy.md` (design contracts, direction). Do not restate them here.

## The rule people trip on

Editing `pkg/rt/core/**/*.lg` changes nothing until `make generate` runs, and the
regenerated files commit together with the edit. `make check-generated` proves freshness;
`make build` does not. Details: `docs/regenerating-generated-artifacts.md`.

## Working here

- `make build`, `make test`, `make lint`, `make generate`; the full map is in
  `docs/contributor-workflow.md`. Install the CI-mirroring hooks once with
  `prek install --install-hooks --hook-type pre-commit --hook-type pre-push`.
- No test files at the repo root; `.lg` tests under `test/`, Go tests beside their package.
- Docs under `docs/` carry frontmatter; run `python3 scripts/docs_frontmatter_hook.py --check`
  on the files you touch.
- Natives register two ways (`ns.Def` in `pkg/rt/lang.go`, or `//lg:native` markers that
  `make generate` turns into registrars). Never both for one name.
- The default compiler is `pkg/compiler`; `*ir-compile*` and `-tags gogen_ir` are the other
  two paths. A change to one needs a stated position on the other two.
- Perf claims go through `docs/perf/ratchet.md`, never a same-machine single run.

## Upstream etiquette (nooga/let-go)

- Changes against a design contract in `docs/contribution-policy.md` need project-owner
  sign-off in the PR thread, not just green CI.
- PR bodies say why before what, describe the net change, and say how it was verified.
  Titles become squash-commit subjects.
- Architecture direction is recorded in issues and in the maintainers' review comments;
  read them before proposing a namespace, IR, or backend change.

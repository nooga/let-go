---
status: active
last-verified: 2026-06-18
human-verified:
---

# Resources and source paths

Two related path-list mechanisms control where `lg` reads files from: resource
roots (non-source data) and source roots (`require`d namespaces). Both take a
path list separated by `:` on Unix, `;` on Windows.

## Resources (`io/resource`)

Programs can read non-source files (templates, static web assets, data) via
`io/resource`, which returns a reader-coercible handle (or `nil` if missing)
that composes with `io/slurp`, `io/reader`, and `io/line-seq`:

```clojure
(when-let [r (io/resource "templates/index.html")]
  (io/slurp r))                     ; => the file contents, or skips if absent
```

Resource roots are given explicitly with `-resource-paths`, or via the
`LG_RESOURCE_PATHS` env var. Resources are addressed by their path relative to a
root; with multiple roots, the first match wins.

```bash
lg -resource-paths resources app.lg          # dev: read from ./resources
```

When you bundle with `-b`, every file under the resource roots is embedded in
the binary, so `io/resource` works on any machine with no files alongside it:

```bash
lg -b myapp -resource-paths resources app.lg  # embed resources into the binary
./myapp                                        # io/resource reads embedded copies
```

A bundled binary reads **only** its embedded resources — it ignores the ambient
filesystem, so deployment is self-contained and predictable. There is no default
resource directory; `lg` is explicit-only.

## Source paths

`require`d namespaces are resolved against a list of search roots. By default
`lg` searches the current directory. You can set the roots explicitly with
`-source-paths` or the `LG_SOURCE_PATHS` env var:

```bash
lg -source-paths src:lib app.lg     # search ./src and ./lib
```

When you provide the search path — by flag or env var — it is taken as the
**complete** list: the current directory is **not** searched implicitly. Add
`.` to the list to include it (`-source-paths .:lib`). A present-but-empty value
(`-source-paths ""` or `LG_SOURCE_PATHS=`) means "no source paths" — only
embedded namespaces resolve. The script passed on the command line is always
loaded by its path, independent of the search path.

If the search path is not given by flag or env var, it defaults to `.` (current
directory).

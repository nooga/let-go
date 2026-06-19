---
status: active
last-verified: 2026-06-19
human-verified:
---

# nREPL

let-go ships an nREPL server that works with CIDER (Emacs), Calva (VS Code), and
Conjure (Neovim).

```bash
lg -n                             # default port 2137
lg -n -p 7888
```

It writes `.nrepl-port` to the working directory so editors auto-discover it.

Supported ops: `clone`, `close`, `eval`, `load-file`, `describe`, `completions`,
`complete`, `info`, `lookup`, `ls-sessions`, `interrupt`.

## Editor setup

- **Emacs (CIDER)**: `M-x cider-connect-clj`, `localhost`, port from `.nrepl-port`.
- **VS Code (Calva)**: open a let-go project (the bundled `.vscode/settings.json`
  registers a connect sequence). Use "Calva: Start a Project REPL and Connect
  (Jack-In)" → "let-go", or "Calva: Connect to a Running REPL Server" if nREPL is
  already up.
- **Neovim (Conjure)**: auto-connects when `.nrepl-port` exists.

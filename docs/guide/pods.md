---
status: active
last-verified: 2026-06-20
human-verified: 2026-06-20
---

# Babashka pods

let-go can load [Babashka pods](https://github.com/babashka/pods), which opens up
the whole pod ecosystem: SQLite, AWS, Docker, file watching, etc.

```clojure
(pods/load-pod 'org.babashka/go-sqlite3 "0.3.13")

(pod.babashka.go-sqlite3/execute! "app.db"
  ["create table users (id integer primary key, name text)"])
(pod.babashka.go-sqlite3/execute! "app.db"
  ["insert into users (name) values ('Alice')"])
(pod.babashka.go-sqlite3/query "app.db"
  ["select * from users"])
;; => [{:id 1 :name "Alice"}]
```

let-go shares `~/.babashka/pods/` with `bb`, so install pods with babashka and
use them from `lg`. See the [pod registry](https://github.com/babashka/pod-registry)
for what's available.

For the host-side protocol and implementation, see
[../design/pods.md](../design/pods.md).

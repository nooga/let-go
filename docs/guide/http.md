---
status: active
last-verified: 2026-09-18
human-verified:
---

# http: serving and fetching over HTTP

The built-in `http` namespace is a Ring-style server plus a small client.
Both are gated off `tinygo` and `lg_no_http` builds.

## Serving

A handler is a function from a request map to a response map:

```clojure
(defn handler [req]
  {:status 200
   :headers {"Content-Type" "text/plain"}
   :body (str "hello " (:uri req))})

(http/serve handler ":8080")   ; blocks until the server stops
```

The request map carries `:request-method` (a lowercase keyword, `:get`),
`:scheme`, `:uri`, `:path`, `:query-string`, `:body` (a string),
`:remote-addr`, `:server-addr`, `:server-port`, `:content-type`, and
`:headers` (lowercased names to comma-joined values).

The response map:

- `:status` — defaults to 200.
- `:headers` — a map of name to value, or absent. An empty map is fine.
- `:body` — a string or a reader is sent buffered. A channel or a lazy
  sequence is streamed, one element per flush, so a handler can pace a
  response (server-sent events, long polls); iteration stops when the
  client goes away.

## Starting and stopping

`http/serve` never returns while the server runs. For a server you can
stop — tests, a REPL, a system that halts its components in order — use
`start`, `stop` and `wait`:

```clojure
(def srv (http/start handler "127.0.0.1:0"))   ; binds now, serves in the background
(:port srv)                                     ; the resolved port
(:addr srv)                                     ; "127.0.0.1:41235"

(http/get (str "http://" (:addr srv) "/x"))

(http/stop srv)        ; graceful: in-flight requests finish, up to 5 s
(http/stop srv 500)    ; then Close after 500 ms
(http/wait srv)        ; blocks until the server has fully stopped; nil
```

`start` returns an `http/Server` record with `:addr`, `:port` and
`:server`, the handle `stop` and `wait` take (either the record or the
handle works). A bind error is thrown by `start` itself. `stop` is
idempotent. `wait` returns `nil` after a clean stop and throws if the
server died on its own.

`serve` is `start` followed by `wait`, so a program's main can also be
written as `(http/wait (http/start handler ":8080"))`.

Both `serve` and `start` run the server under the caller's scope:
cancelling that scope stops the server the way `stop` does.

## Client

```clojure
(http/get "https://example.com/")                        ; → {:status :headers :body}
(http/post url "payload" {:headers {"Content-Type" "text/plain"}})
(http/request {:method :put :url url :body "..." :headers {...}})
(http/get url {:as :stream})                             ; :body is an io reader
```

Responses are maps with `:status`, `:headers` (lowercased names) and
`:body`. Requests inherit the caller's scope, so a cancelled scope aborts
them.

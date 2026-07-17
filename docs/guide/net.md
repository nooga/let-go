---
status: active
last-verified: 2026-07-17
human-verified:
---

# net and bencode: TCP clients in pure let-go

The `net` namespace provides a minimal TCP client, and `bencode` provides
bencode framing over a connection — together enough to write network tools
(an nREPL client, for one) in pure let-go. Both namespaces are gated off
`js`/`wasm` builds, where sockets don't exist.

## net

```clojure
(def conn (net/dial "localhost" 7888)) ; connect, 3s timeout
(net/write! conn "hello")              ; string or byte-array, raw bytes on the wire
(net/read! conn 4096)                  ; → byte-array, nil on clean EOF
(net/close! conn)
```

`net/dial` takes a host string and an integer port and returns a connection
value. `net/read!` blocks until data arrives, returning at most `max-bytes`
bytes as a byte-array, or `nil` when the peer closes cleanly.

## bencode

`bencode/write!` encodes one value onto a connection: maps (string or keyword
keys), sequences, integers, and strings. `bencode/read!` blocks for the next
value and returns `nil` on clean EOF; the two-arg form takes a timeout in
milliseconds and errors when it elapses.

```clojure
(bencode/write! conn {:op "eval" :code "(+ 1 2)"})
(bencode/read! conn)      ; blocks for the response
(bencode/read! conn 2000) ; same, but errors after 2s
```

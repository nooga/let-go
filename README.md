<!--suppress ALL -->
<p align="center">
<img src="meta/logo.png" alt="Squishy loafer" title="Squishy loafer of let-go" />
</p>


![Tests](https://github.com/nooga/let-go/actions/workflows/go.yml/badge.svg)

# let-go

Greetings loafers! *(λ-gophers haha, get it?)*

This is a bytecode compiler and VM for a language closely resembling Clojure, a Clojure dialect, if you will. 

Here are some nebulous goals in no particular order:
- [x] Quality entertainment,
- [ ] Implement as much of Clojure as possible - including persistent data types and true concurrency,
- [x] Provide comfy two-way interop for arbitrary functions and types,
- [ ] Serve primarily as an embedded extension language,
- [ ] AOT (let-go -> standalone binary) would be nice eventually, 
- [ ] Strech goal: let-go bytecode -> Go translation.
- [ ] ~~Pure Go, zero dependencies.~~ We use a lightweight line editor for the REPL now - [alimpfard/line](https://github.com/alimpfard/line)

Here are the non goals:
- Stellar performance (btw. it seems to be faster than Joker),
- Being a drop-in replacement for [clojure/clojure](https://github.com/clojure/clojure) at any point,
- Being a linter/formatter/tooling for Clojure in general.

## Current status 

It can eval basic Clojure-like lisp. It's semi-good for solving Advent of Code problems but not for anything more serious yet. 

Everything more or less half-baked and most features are happy path only ;)

## Examples

See [tests](https://github.com/nooga/let-go/tree/main/test) and [examples](https://github.com/nooga/let-go/tree/main/examples) for some code samples. 

## Try online

Check out [this bare-bones online REPL](https://nooga.github.io/let-go/).

## Prerequisites and installation

Building or running let-go from source requires Go 1.19. 

```
go install github.com/nooga/let-go@latest
```

## Building

The best way to play with `let-go` right now is to clone this repo and run the REPL like this:

```
go run . 
```

To run an expression:

```
go run . -e '(+ 1 1)'
```

To run a file:

```
go run . test/hello.lg
```

Use the `-r` flag to run the REPL after the interpreter has finished with files and `-e`:

```bash
go run . -r test/simple.lg                # will run simple.lg first, then open up a REPL
go run . -r -e '(* fun 2)' test/simple.lg # will run simple.lg first, then (* fun 2) and REPL 
```

## Building the interpreter -`lg`

To build and start the standalone interpreter:
```bash
make
```

---
Follow me on twitter for nightly updates! 🌙

<a href="https://twitter.com/intent/follow?screen_name=mgasperowicz">
<img src="https://img.shields.io/twitter/follow/mgasperowicz?style=social&logo=twitter"
alt="follow on Twitter"></a>

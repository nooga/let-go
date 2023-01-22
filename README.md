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
- [ ] Making it legal to write Clojure at your Go dayjob,
- [ ] Implement as much of Clojure as possible - including persistent data types and true concurrency,
- [x] Provide comfy two-way interop for arbitrary functions and types,
- [ ] AOT (let-go -> standalone binary) would be nice eventually, 
- [ ] Strech goal: let-go bytecode -> Go translation.

Here are the non goals:
- Stellar performance (cough cough, it seems to be faster than [Joker](https://github.com/candid82/joker)),
- Being a drop-in replacement for [clojure/clojure](https://github.com/clojure/clojure) at any point,
- Being a linter/formatter/tooling for Clojure in general.

## Current status 

It can eval a Clojure-like lisp. It's semi-good for [solving Advent of Code problems](https://github.com/nooga/aoc2022) but not for anything more serious yet. 

There are some goodies:

- [x] Macros with syntax quote,
- [x] Reader conditionals,
- [x] Destructuring,
- [x] Multi-arity functions,
- [x] Atoms, channels & go-blocks a'la `core.async`, 
- [x] Regular expressions (the Go flavor),
- [x] Simple `json`, `http` and `os` namespaces,
- [x] Many functions ported from `clojure.core`,
- [x] REPL with syntax-highlighting and completions,
- [x] Simple nREPL server that seems to work with [BetterThanTomorrow/calva](https://github.com/BetterThanTomorrow/calva),

I'm currently contemplating:

- [ ] Dynamic variables,
- [ ] Exceptions & more comprehensive runtime errors,
- [ ] Numeric tower i.e. floats, bignums and other exotic number species,
- [ ] Optimized persistent data-structures and laziness,
- [ ] Records and protocols,
- [ ] A real test suite would be nice,

## Examples

See:
- [My AoC 2022 solutions](https://github.com/nooga/aoc2022) for an idea of how let-go looks like,
- [Examples](https://github.com/nooga/let-go/tree/main/examples) for small programs I wrote on a whim,
- [Tests](https://github.com/nooga/let-go/tree/main/test) for some random "tests".


## Try online

Check out [this bare-bones online REPL](https://nooga.github.io/let-go/). It runs a WASM build of let-go in your browser! 

## Prerequisites and installation

Building or running let-go from source requires Go 1.19. 

```
go install github.com/nooga/let-go@latest
```

Try it out:

```
let-go
```

## Running from source

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

---
[🤓 Follow me on twitter](https://twitter.com/MGasperowicz)
[🐬 Check out monk.io](https://monk.io)

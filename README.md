<!--suppress ALL -->
<p align="center">
<img src="meta/logo.png" alt="Squishy loafer" title="Squishy loafer of let-go" />
</p>


![Tests](https://github.com/nooga/let-go/actions/workflows/go.yml/badge.svg)

# let-go

Greetings loafers! *(λ-gophers haha, get it?)*

This is supposed to be a compiler and bytecode VM for a language resembling Clojure as close as possible.

Now, I know about [candid82/joker](https://github.com/candid82/joker) and I 💛 it. Though, it has some 
drawbacks and design choices that I'd like to avoid.

Here are some nebulous goals in no particular order:
- Quality entertainment,
- Implement as much of Clojure as possible - including persistent data types and true concurrency,
- Provide comfy two-way interop for arbitrary functions and types,
- Serve primarily as an embedded extension language,
- Standalone interpreter mode and AOT (let-go -> standalone binary) would be nice eventually, 
- Strech goal: let-go bytecode -> Go translation. 

Here are the non goals:
- Stellar performance,
- Being a drop-in replacement for [clojure/clojure](https://github.com/clojure/clojure) at any point,
- Being a linter/formatter/tooling for Clojure in general.

## Current status 

It doesn't do much yet - just started laying down the basics. Come back in a couple of months to see cool demos.

See [compiler tests](https://github.com/nooga/let-go/blob/master/pkg/compiler/compiler_test.go) if you're really interested. 
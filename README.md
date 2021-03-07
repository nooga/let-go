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
- Pure Go, zero dependencies.

Here are the non goals:
- Stellar performance,
- Being a drop-in replacement for [clojure/clojure](https://github.com/clojure/clojure) at any point,
- Being a linter/formatter/tooling for Clojure in general.

## Current status 

Can compile and eval basic Clojure flavored lisp.

#### Notable features
- ✅ macros
- ✅ vars
- ✅ can compile and run Y Combinator

#### The most impressive snippet so far

```clojure
(test "Y combinator"
      (let [Y (fn [f] ((fn [x] (x x))
                        (fn [x]
                          (f (fn [y] ((x x) y))))))
            fac-gen (fn [func] (fn [n] (if (zero? n) 1 (* n (func (dec n))))))]
        (= 120 ((Y fac-gen) 5))))
```

See [tests](https://github.com/nooga/let-go/tree/main/test) for more examples. 

## Running

Sure, you can! Just keep in mind that we're not there yet and it will likely blow up in your 
face. Just remember to leave an issue when it does 😊

The best way to take `let-go` for a spin right now is to clone this repo and run the REPL like this:

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
Follow me on twitter for nightly updates! 🌙

<a href="https://twitter.com/intent/follow?screen_name=mgasperowicz">
<img src="https://img.shields.io/twitter/follow/mgasperowicz?style=social&logo=twitter"
alt="follow on Twitter"></a>

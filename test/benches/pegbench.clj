;; PEG-combinator-shaped microbenchmark.
;; Mirrors yamlstar's hot pattern: closures built by combinators,
;; invoked indirectly through vars, threading (string, pos) state.
;; Portable across Clojure / gloat / let-go: no interop, no metadata.

(defn chr [c]
  (fn [s pos]
    (if (and (< pos (count s)) (= (subs s pos (inc pos)) c))
      (inc pos)
      -1)))

(defn alt2 [a b]
  (fn [s pos]
    (let [r (a s pos)]
      (if (neg? r) (b s pos) r))))

(defn cat2 [a b]
  (fn [s pos]
    (let [r (a s pos)]
      (if (neg? r) -1 (b s r)))))

(defn rep* [p]
  (fn [s pos]
    (loop [pos pos]
      (let [r (p s pos)]
        (if (neg? r) pos (recur r))))))

;; grammar: ( ('a'|'b') 'c' )*
;; Built inside a fn (not top-level def) because glojure AOT cannot
;; serialize closures computed at load time ("opaque go function values").
(defn make-grammar []
  (rep* (cat2 (alt2 (chr "a") (chr "b")) (chr "c"))))

(defn run-bench [n]
  (let [grammar (make-grammar)
        s "acbcacbcacbcacbcacbcacbcacbcacbc"]
    (loop [i 0 acc 0]
      (if (< i n)
        (recur (inc i) (+ acc (grammar s 0)))
        acc))))

(defn -main [& args]
  (println (run-bench 200000)))

;; Gold differential test: dynamic bindings must convey through EAGER
;; callback-invoking builtins (mapv, reduce, swap!, transducers, …) exactly as
;; in Clojure — including when the binding lives in a child ExecContext (inside
;; a future), where a context-free Fn.Invoke would fall back to the root.
;;
;; Run the SAME file under both runtimes and diff stdout:
;;   clojure -M dynvar_callbacks.cljc
;;   lg          dynvar_callbacks.cljc
;; The single printed line must match byte-for-byte.
;;
;; The eager/lazy split is intentional: eager callbacks run synchronously
;; inside the binding's extent, so they MUST see the binding. Lazy seqs
;; (map/filter) realized after the extent are NOT exercised here — Clojure
;; itself documents that caveat, and let-go matches it.

(def ^:dynamic *x* 1)

(defn scenarios []
  (binding [*x* 10]
    @(future
       (binding [*x* 100]
         [;; 1. mapv: eager map, callback must see the future's binding.
          (mapv (fn [_] *x*) [0 0])

          ;; 2. filterv: eager filter via predicate callback.
          (filterv (fn [_] (= *x* 100)) [:keep])

          ;; 3. reduce: accumulator callback.
          (reduce (fn [a _] (+ a *x*)) 0 [1 1])

          ;; 4. run!: side-effecting eager traversal.
          (let [acc (atom [])]
            (run! (fn [_] (swap! acc conj *x*)) [0])
            @acc)

          ;; 5. every?: short-circuiting eager predicate.
          (every? (fn [_] (= *x* 100)) [0])

          ;; 6. sort-by: key-fn callback. Ascending iff the binding is seen.
          (sort-by (fn [n] (if (= *x* 100) n (- n))) [2 1 3])

          ;; 7. swap!: atom update fn.
          (let [a (atom 0)] (swap! a (fn [_] *x*)) @a)

          ;; 8. transducer application: into with an xform.
          (into [] (map (fn [_] *x*)) [0])

          ;; 9. with-out-str: *out* rebound in the child context, print
          ;;    must write to the binding's writer, not the root stdout.
          (with-out-str (print *x*))]))))

(println (scenarios))

;; No-op in let-go; releases Clojure's future thread-pool so the JVM exits
;; promptly instead of lingering ~60s on the agent keep-alive (this script is
;; only run under real Clojure when re-deriving the gold via LETGO_GOLD_REDERIVE).
(shutdown-agents)

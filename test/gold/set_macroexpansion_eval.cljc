;; JVM-Clojure oracle for set-literal evaluation across macroexpansion.
;;
;; The two source forms below are distinct, but the one-step threading macro
;; expands to the same (swap! calls inc) form as its sibling. Clojure evaluates
;; both elements before constructing/deduplicating the runtime set. A compiler
;; must not rebuild the expanded ASTs as a host-language set and lose one
;; evaluation.

(ns test.gold.set-macroexpansion-eval)

(defn run []
  (let [calls  (atom 0)
        values #{(-> (swap! calls inc))
                 (swap! calls inc)}]
    [@calls (count values) (contains? values 1) (contains? values 2)]))

(println (pr-str (run)))
(shutdown-agents)

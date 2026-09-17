;; JVM oracle: macroexpansion must not collapse distinct map entry slots.
(ns test.gold.map-macroexpansion-eval)

(defn run []
  (let [calls (atom 0)
        values {(-> (swap! calls inc)) (swap! calls inc)
                (swap! calls inc)      (swap! calls inc)}]
    [@calls
     (count values)
     (get values 1)
     (get values 3)
     (contains? values 1)
     (contains? values 3)]))

(println (pr-str (run)))
(shutdown-agents)

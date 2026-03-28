;; Build a persistent map with 10000 entries
(reduce (fn [m i] (assoc m i (* i i)))
        {}
        (range 10000))

;; Transducer pipeline — no intermediate collections
(transduce
  (comp (map #(* % %))
        (filter even?)
        (take 100))
  + 0
  (range 10000))

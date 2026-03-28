;; Map + filter + take pipeline over lazy seqs
(reduce + 0
  (take 100
    (filter even?
      (map #(* % %) (range 10000)))))

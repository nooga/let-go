(defn tak [x y z]
  (if (< y x)
    (tak (tak (dec x) y z) (tak (dec y) z x) (tak (dec z) x y))
    z))

(println (tak 30 22 12))

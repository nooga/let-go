(ns test.gold.return-hint-metadata)

(def return-hinted-source "(defn hinted ^long [^long n] n)")

(defn reader-metadata []
  (let [arity (nth (read-string return-hinted-source) 2)]
    ;; JVM Clojure attaches reader metadata directly. let-go preserves the
    ;; reader form as (with-meta value metadata), so normalize both shapes to
    ;; the same metadata maps before comparing their behavior.
    (if (vector? arity)
      [(meta arity) (meta (first arity))]
      (let [arg        (first (second arity))
            arity-meta (nth arity 2)
            arg-meta   (nth arg 2)]
        [(assoc arity-meta :tag (second (:tag arity-meta)))
         (assoc arg-meta :tag (second (:tag arg-meta)))]))))

(defn hinted ^long [^long n] n)

(defn run ^Object []
  [(hinted 7) (reader-metadata)])

(println (pr-str (run)))

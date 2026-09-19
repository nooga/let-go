(ns test.gold.return-hint-metadata)

(def return-hinted-source "(defn hinted ^long [^long n] n)")

(defn reader-metadata []
  (let [arity (nth (read-string return-hinted-source) 2)]
    ;; JVM Clojure attaches reader metadata directly. let-go's read-string
    ;; attaches it to the vector too, but a let-go symbol cannot carry
    ;; metadata, so the hinted parameter still reads as the form
    ;; (with-meta n {:tag (quote long)}); an older code-mode reader returned
    ;; that form for the vector as well. Normalize every shape to the same
    ;; metadata maps before comparing their behavior.
    (if (vector? arity)
      (let [param      (first arity)
            param-meta (if (symbol? param)
                         (meta param)
                         (let [m (nth param 2)]
                           (assoc m :tag (second (:tag m)))))]
        [(meta arity) param-meta])
      (let [arg        (first (second arity))
            arity-meta (nth arity 2)
            arg-meta   (nth arg 2)]
        [(assoc arity-meta :tag (second (:tag arity-meta)))
         (assoc arg-meta :tag (second (:tag arg-meta)))]))))

(defn hinted ^long [^long n] n)

(defn run ^Object []
  [(hinted 7) (reader-metadata)])

(println (pr-str (run)))

;; JVM-Clojure oracle for what a macro sees when its argument is a map or set
;; literal that references a deftype or defrecord field.
;;
;; A field is a lexical local inside the method body, so the macro receives the
;; literal exactly as written: a map is still map?, a set is still set?. A
;; compiler that rewrites field references before expanding macros hands the
;; macro a constructor call instead, and the macro takes its :other branch.

(ns test.gold.macro-literal-kind-field)

(defmacro literal-kind [form]
  (cond (map? form) :map (set? form) :set :else :other))

(defprotocol P (probe [this]))

(deftype T [x]
  P
  (probe [this]
    [(literal-kind {x 1}) (literal-kind {:k x}) (literal-kind #{x})
     (literal-kind [x]) (when true (literal-kind {x x}))]))

(defrecord R [x]
  P
  (probe [this]
    [(literal-kind {x 1}) (literal-kind {:k x}) (literal-kind #{x})
     (literal-kind [x]) (when true (literal-kind {x x}))]))

(println (pr-str [(probe (->T :field)) (probe (->R :field))]))
(shutdown-agents)

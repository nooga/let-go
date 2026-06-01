;; Proof that a require under a feature let-go does NOT match is skipped in
;; its entirety. The namespace below does not exist; if the :cljs splice were
;; not fully dropped at read time, loading this file would fail with a
;; "could not locate" error. This is the exact mechanism by which JVM Clojure
;; skips :lg requires — so a Clojure consumer of shared code never attempts to
;; load a let-go-only namespace guarded behind #?@(:lg ...).
(ns portable-foreign
  #?(:cljs
     (:require
      [totally.bogus.absent :as bogus])))

(def loaded :ok)

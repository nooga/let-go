;; A portable .cljc consumer. The let-go-specific dependency (let-go.semver)
;; and every use of it are guarded behind :lg reader conditionals, so a
;; JVM-Clojure reader — which never matches :lg — skips them entirely and
;; never tries to load the let-go-only namespace. Under let-go the :lg
;; branches are taken and the semver capability is used.
(ns portable-consumer
  #?(:lg
     (:require
      [let-go.semver :as semver])))

(defn version-render
  "Under let-go, normalize a version string via let-go.semver; on any other
   host the :lg branch is skipped and the string passes through. Note the
   semver alias appears ONLY inside the :lg branch, so a non-let-go reader
   never sees an unresolved symbol."
  [s]
  #?(:lg (semver/render (semver/version s))
     :default s))

(defn at-least?
  "Under let-go, true when v >= floor by semver precedence; elsewhere a
   conservative pass-through."
  [v floor]
  #?(:lg (semver/gte? v floor)
     :default true))

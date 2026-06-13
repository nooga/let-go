;; Gold differential test: the *ns* dynamic var (current namespace) must isolate
;; and convey across nested threads exactly like Clojure. Run the same file
;; under both runtimes and diff stdout (see run_gold.sh); the printed line must
;; match byte-for-byte.
;;
;; KNOWN LIMITATION (intentionally NOT exercised here): let-go's compiler detects
;; a leading top-level (in-ns 'x) at compile time (tryDetectInNS) and switches
;; the namespace via the SHARED *ns* root, which leaks to runtime. So
;;   (binding [*ns* a] (in-ns 'b) ...)
;; — with in-ns as the FIRST body form — diverges from Clojure (it mutates the
;; root). Scenario 4 places a no-op before in-ns so the compile-time detector
;; does not fire; the runtime in-ns then correctly mutates only the thread-local
;; binding. The compile-time conflation is a separate compiler concern.

(create-ns 'gold.a)
(create-ns 'gold.b)

(defn scenarios []
  [;; 1. *ns* conveys into a future (binding conveyance at creation).
   (ns-name (binding [*ns* (find-ns 'gold.a)] @(future *ns*)))

   ;; 2. Concurrent futures each rebinding *ns* do not cross-talk, and the
   ;;    parent's *ns* is restored afterwards.
   (let [fa (future (binding [*ns* (find-ns 'gold.a)] (ns-name *ns*)))
         fb (future (binding [*ns* (find-ns 'gold.b)] (ns-name *ns*)))]
     [@fa @fb (ns-name *ns*)])

   ;; 3. A future nested inside a future inherits the ENCLOSING future's *ns*.
   (ns-name (binding [*ns* (find-ns 'gold.a)]
              @(future (binding [*ns* (find-ns 'gold.b)] @(future *ns*)))))

   ;; 4. in-ns mutates the thread-local *ns* (thread-local set!), isolated to
   ;;    the future; the parent root is untouched. (in-ns is not the first body
   ;;    form — see the KNOWN LIMITATION note above.)
   (let [f (future (binding [*ns* (find-ns 'gold.a)] :nop (in-ns 'gold.b) (ns-name *ns*)))]
     [@f (ns-name *ns*)])])

(println (pr-str (scenarios)))
(shutdown-agents)

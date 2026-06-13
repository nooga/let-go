;; Gold differential test: dynamic-var (namespace variable) binding semantics
;; across nested threads must be IDENTICAL between Clojure and let-go.
;;
;; Run the SAME file under both runtimes and diff stdout:
;;   clojure -M dynvar_threads.cljc
;;   lg          dynvar_threads.cljc
;; (see run_gold.sh). The single printed line must match byte-for-byte.
;;
;; This exercises the ExecContext binding-propagation model: a future captures
;; the dynamic bindings active where it is created (Clojure "binding
;; conveyance"), is isolated from later rebinds on other threads, and a future
;; nested inside a future inherits the ENCLOSING future's bindings.
;;
;; PERMISSIVE DEVIATIONS (intentionally NOT exercised here, so the script stays
;; in the subset where let-go and Clojure agree byte-for-byte):
;;   a) (set! *v* val) with NO active binding: Clojure throws "Can't set!: *v*
;;      from non-binding thread"; let-go permits it and mutates the root.
;;   b) (set! *v* val) on a CONVEYED binding inside a future (a future that did
;;      not itself open a `binding`): Clojure throws "non-binding thread";
;;      let-go permits it (the conveyed frame is a real binding in the child
;;      ExecContext). Every future below opens its own `binding` before set!,
;;      which is well-defined in both runtimes.

(def ^:dynamic *x* :root)

;; A protocol method and a multimethod that both read the dynamic var. Their
;; dispatch must thread the ExecContext so the method body sees the caller's
;; (here: the future's) binding, not the process root.
(defprotocol PRead (pread [x]))
(deftype TRead [] PRead (pread [_] *x*))
(defmulti mread (fn [_] :k))
(defmethod mread :k [_] *x*)

(defn scenarios []
  [;; 1. A future inherits the caller's binding (conveyance at creation time).
   (binding [*x* :a] @(future *x*))

   ;; 2. A future captures its creation-time binding; a later rebind on the
   ;;    CALLING thread does not change what the already-spawned future sees.
   (binding [*x* :b]
     (let [f (future *x*)]
       (binding [*x* :changed] @f)))

   ;; 3. A future nested inside a future inherits the ENCLOSING future's
   ;;    binding (:inner), not the process root (:root) or the outer (:outer).
   (binding [*x* :outer]
     @(future (binding [*x* :inner] @(future *x*))))

   ;; 4. Concurrent futures with distinct rebindings do not cross-talk; each
   ;;    sees only its own. Derefed in a fixed order so output is deterministic.
   (binding [*x* :base]
     (let [fa (future (binding [*x* :p] *x*))
           fb (future (binding [*x* :q] *x*))]
       [@fa @fb]))

   ;; 5. bound-fn closes over the dynamic scope active where it was created and
   ;;    re-establishes it on each later call, even outside that binding.
   (let [g (binding [*x* :captured] (bound-fn [] *x*))]
     (g))

   ;; --- mutation (set!) ---------------------------------------------------
   ;; 6. set! mutates the CURRENT thread-local binding: visible to later reads
   ;;    within the binding, and it does NOT leak past the binding (scenario 9
   ;;    confirms the root is still :root afterwards).
   (binding [*x* :a] (set! *x* :m) *x*)

   ;; 7. A future that establishes its OWN binding may set! within it; the
   ;;    mutation is isolated to that future and never touches the parent's
   ;;    binding.
   (binding [*x* :base]
     (let [f (future (binding [*x* :f] (set! *x* :fm) *x*))]
       [@f *x*]))

   ;; 8. Concurrent futures each mutating their own binding do not cross-talk.
   (binding [*x* :base]
     (let [fa (future (binding [*x* :pa] (set! *x* :ma) *x*))
           fb (future (binding [*x* :pb] (set! *x* :mb) *x*))]
       [@fa @fb *x*]))

   ;; 9. Protocol-method and multimethod dispatch convey the binding: invoked
   ;;    inside a future that established its own binding, both method bodies
   ;;    read the future's *x* (:pm), not the process root (:root). (Regression:
   ;;    ProtocolFn/MultiFn dispatch must route the impl through the ExecContext.)
   (binding [*x* :base]
     @(future (binding [*x* :pm] [(pread (->TRead)) (mread 1)])))

   ;; 10. The parent root binding is fully restored after all the threaded work
   ;;     AND all the set! mutations above — none of them leaked to the root.
   *x*])

(println (pr-str (scenarios)))

;; No-op in let-go; releases Clojure's future thread-pool so the JVM exits
;; promptly instead of lingering ~60s on the agent keep-alive.
(shutdown-agents)

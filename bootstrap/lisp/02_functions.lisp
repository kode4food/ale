;;;; ale bootstrap: functions

(defmacro assert-args
  ([] nil)
  ([clause]
   (raise "assert-args clauses must be paired"))
  ([& clauses]
   `(if ~(clauses 0)
      (assert-args ~@(rest (rest clauses)))
      (raise ~(clauses 1)))))

(defn partial
  ([func] func)
  ([func & first-args]
   (assert-args
    (is-apply func) "partial requires a function")
   (fn [& rest-args]
     (apply func (append first-args rest-args)))))

(defn no-op
  [& _])

(defn identity
  [val]
  val)

(defn constantly
  [val]
  (fn [& _] val))

(defn comp-outer
  [func args rest-funcs]
  (if (seq rest-funcs)
    (comp-outer (first rest-funcs) (list func args) (rest rest-funcs))
    (list func args)))

(defmacro comp
  ([] identity)
  ([func] func)
  ([func & funcs]
   (let [args        (gensym "args")
         inner       (list 'apply func args)
         first-outer (first funcs)
         rest-outer  (rest funcs)]
     `(fn [& ~args]
        ~(comp-outer first-outer inner rest-outer)))))

(defmacro juxt
  [& funcs]
  (let [args (gensym "args")]
    `(fn [& ~args]
       [~@(map (fn [f] (list 'apply f args)) funcs)])))

(defmacro .
  [target method & args]
  `((get ~target ~method) ~@args))

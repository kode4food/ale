;;;; ale core: functions

(defmacro assert-args
  ([] nil)
  ([clause]
    (raise "assert-args clauses must be paired"))
  ([& clauses]
    `(if ~(clauses 0)
         (assert-args ~@(rest (rest clauses)))
         (raise ~(clauses 1)))))

(defn no-op
  [& _])

(defn identity
  [val]
  val)

(defn constantly
  [val]
  (fn [& _] val))

(defmacro .
  [target method & args]
  `((get ~target ~method) ~@args))

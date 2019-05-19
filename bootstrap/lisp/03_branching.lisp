;;;; ale bootstrap: branching

(defmacro not
  [val]
  `(if ~val false true))

(defmacro unless
  ([test] nil)
  ([test then] `(if ~test nil ~then))
  ([test then else] `(if ~test ~else ~then)))

(defmacro when
  ([test] nil)
  ([test form] `(if ~test ~form nil))
  ([test & forms] `(if ~test (do ~@forms) nil)))

(defmacro when-not
  ([test] nil)
  ([test form] `(if ~test nil ~form))
  ([test & forms] `(if ~test nil (do ~@forms))))

(defmacro and
  ([] true)
  ([clause] clause)
  ([& clauses]
   `(let [and# ~(clauses 0)]
      (if and# (and ~@(rest clauses)) and#))))

(defmacro !and
  [& clauses]
  `(not (and ~@clauses)))

(defmacro or
  ([] nil)
  ([clause] clause)
  ([& clauses]
   `(let [or# ~(clauses 0)]
      (if or# or# (or ~@(rest clauses))))))

(defmacro !or
  [& clauses]
  `(not (or ~@clauses)))

(defmacro cond
  ([] nil)
  ([clause] clause)
  ([& clauses]
   (let [test   (clauses 0)
         branch (clauses 1)]
     (unless (and (is-atom test) test)
             `(if ~test
                ~branch
                (cond ~@(rest (rest clauses))))
             branch))))

(defmacro if-let
  ([binding then] `(if-let ~binding ~then nil))
  ([binding then else]
   (assert-args
    (is-vector binding) "binding vector must be supplied"
    (= 2 (size binding)) "binding vector must contain 2 elements")
   (let [sym  (binding 0)
         test (binding 1)]
     `(let [~sym ~test]
        (if ~sym ~then ~else)))))

(defmacro when-let
  ([binding form] `(if-let ~binding ~form))
  ([binding & body] `(if-let ~binding (do ~@body))))

;;;; ale core: branching

(define not is-false)

(defmacro unless
  ([test] null)
  ([test then]      `(if ,test null ,then))
  ([test then else] `(if ,test ,else ,then)))

(defmacro when
  ([test] null)
  ([test form]    `(if ,test ,form null))
  ([test . forms] `(if ,test (do ,@forms) null)))

(defmacro when-not
  ([test] null)
  ([test form]    `(if ,test null ,form))
  ([test . forms] `(if ,test null (do ,@forms))))

(defmacro and
  ([] true)
  ([clause] clause)
  (clauses
    `(let [and# ,(clauses 0)]
       (if and#
           (and ,@(rest clauses))
           and#))))

(defmacro !and clauses
  `(not (and ,@clauses)))

(defmacro or
  ([] null)
  ([clause] clause)
  (clauses
    `(let [or# ,(clauses 0)]
       (if or#
           or#
           (or ,@(rest clauses))))))

(defmacro !or clauses
  `(not (or ,@clauses)))

(defmacro cond
  ([] null)
  ([clause] clause)
  (clauses
    (let [test   (clauses 0)
          branch (clauses 1)]
      (unless (and (is-atom test) test)
              `(if ,test
                   ,branch
                   (cond ,@(rest (rest clauses))))
              branch))))

(defmacro if-let
  ([binding then] `(if-let ,binding ,then null))
  ([binding then else]
    (assert-args
      (is-vector binding)    "binding vector must be supplied"
      (= 2 (length binding)) "binding vector must contain 2 elements")
    (let [sym  (binding 0)
          test (binding 1)]
      `(let [,sym ,test]
            (if ,sym ,then ,else)))))

(defmacro when-let
  ([binding form]   `(if-let ,binding ,form))
  ([binding . body] `(if-let ,binding (do ,@body))))

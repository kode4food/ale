;;;; ale core: branching

(define not is-false)

(define-macro unless
  (lambda
    ([test] '())
    ([test then]      `(if ,test '() ,then))
    ([test then else] `(if ,test ,else ,then))))

(define-macro when
  (lambda
    ([test] '())
    ([test form]    `(if ,test ,form '()))
    ([test . forms] `(if ,test (do ,@forms) '()))))

(define-macro when-not
  (lambda
    ([test] '())
    ([test form]    `(if ,test '() ,form))
    ([test . forms] `(if ,test '() (do ,@forms)))))

(define-macro and
  (lambda
    ([] #t)
    ([clause] clause)
    (clauses
      `(let [and# ,(clauses 0)]
        (if and#
            (and ,@(rest clauses))
            and#)))))

(define-macro (!and . clauses)
  `(not (and ,@clauses)))

(define-macro or
  (lambda
    ([] '())
    ([clause] clause)
    (clauses
      `(let [or# ,(clauses 0)]
        (if or#
            or#
            (or ,@(rest clauses)))))))

(define-macro (!or . clauses)
  `(not (or ,@clauses)))

(define-macro cond
  (lambda
    ([] '())
    ([clause] clause)
    (clauses
      (let [test   (clauses 0)
            branch (clauses 1)]
        (unless (and (is-atom test) test)
                `(if ,test
                    ,branch
                    (cond ,@(rest (rest clauses))))
                branch)))))

(define-macro if-let
  (lambda
    ([binding then] `(if-let ,binding ,then '()))
    ([binding then else]
      (assert-args
        (is-vector binding)    "binding vector must be supplied"
        (= 2 (length binding)) "binding vector must contain 2 elements")
      (let [sym  (binding 0)
            test (binding 1)]
        `(let [,sym ,test]
              (if ,sym ,then ,else))))))

(define-macro when-let
  (lambda
    ([binding form]   `(if-let ,binding ,form))
    ([binding . body] `(if-let ,binding (do ,@body)))))

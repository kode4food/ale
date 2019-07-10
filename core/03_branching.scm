;;;; ale core: branching

(define not is-false)

(define-macro unless
  (lambda
    [(test)           '()]
    [(test then)      `(if ,test '() ,then)]
    [(test then else) `(if ,test ,else ,then)]))

(define-macro when
  (lambda
    [(test)         '()]
    [(test form)    `(if ,test ,form '())]
    [(test . forms) `(if ,test (begin ,@forms) '())]))

(define-macro when-not
  (lambda
    [(test)         '()]
    [(test form)    `(if ,test '() ,form)]
    [(test . forms) `(if ,test '() (begin ,@forms))]))

(define-macro and
  (lambda
    [() #t]
    [(clause) clause]
    [clauses
      `(let [and# ,(clauses 0)]
        (if and#
            (and ,@(rest clauses))
            and#))]))

(define-macro (!and . clauses)
  `(not (and ,@clauses)))

(define-macro or
  (lambda
    [() '()]
    [(clause) clause]
    [clauses
      `(let [or# ,(clauses 0)]
        (if or#
            or#
            (or ,@(rest clauses))))]))

(define-macro (!or . clauses)
  `(not (or ,@clauses)))

(define-macro cond
  (lambda
    [() '()]
    [clauses
      (let [clause (first clauses)]
        (assert-args
          (is-vector clause)    "cond clauses must be vectors"
          (= 2 (length clause)) "cond clauses must be paired")
        (let [test   (clause 0)
              branch (clause 1)]
          `(if ,test
               ,branch
               (cond ,@(rest clauses)))))]))

;; case requires it
(define (map! func coll)
  (unless (is-empty coll)
          (cons (func (first coll)) (map! func (rest coll)))
          '()))

(define-macro (case expr . cases)
  (letrec [val       (gensym "val")
           pred-list (lambda (l) `(or ,@(map! pred l)))
           pred      (lambda (x) `(eq ,val ,x))

           case*
           (lambda
             [() '(raise "no cases could be matched")]
             [clauses
               (let [clause (first clauses)]
                 (assert-args
                   (is-vector clause)    "case clauses must be vectors"
                   (= 2 (length clause)) "case clauses must be paired")
                 (let [test   (clause 0)
                       branch (clause 1)
                       next   (rest clauses)]
                   `(if ,(if (is-list test)
                             (pred-list test)
                             (pred test))
                        ,branch
                        ,(apply case* next))))])]
    `(let [,val ,expr]
        ,(apply case* cases))))

(define-macro if-let
  (lambda
    [(binding then) `(if-let ,binding ,then '())]
    [(binding then else)
      (assert-args
        (is-vector binding)    "binding vector must be supplied"
        (= 2 (length binding)) "binding vector must contain 2 elements")
      (let [sym  (binding 0)
            test (binding 1)]
        `(let [,sym ,test]
              (if ,sym ,then ,else)))]))

(define-macro when-let
  (lambda
    [(binding form)   `(if-let ,binding ,form)]
    [(binding . body) `(if-let ,binding (begin ,@body))]))

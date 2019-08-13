;;;; ale core: branching

(define not is-false)

(define-macro unless
  [(test)           '()]
  [(test then)      `(if ,test '() ,then)]
  [(test then else) `(if ,test ,else ,then)])

(define-macro when
  [(test)         '()]
  [(test form)    `(if ,test ,form '())]
  [(test . forms) `(if ,test (begin ,@forms) '())])

(define-macro when-not
  [(test)         '()]
  [(test form)    `(if ,test '() ,form)]
  [(test . forms) `(if ,test '() (begin ,@forms))])

(define-macro cond
  [() '()]
  [clauses
    (let [clause (first clauses)]
      (assert-args
        (is-vector clause)    (str "cond clause must be a vector: " clause)
        (= 2 (length clause)) (str "cond clause must be paired: " clause))
      (let ([test   (clause 0)]
            [branch (clause 1)])
        `(if ,test
              ,branch
              (cond ,@(rest clauses)))))])

;; case requires it
(define (map! func coll)
  (unless (is-empty coll)
          (cons (func (first coll)) (map! func (rest coll)))
          '()))

(define-macro (case expr . cases)
  (let-rec ([val       (gensym "val")]
            [pred-list (lambda (l) `(or ,@(map! pred l)))]
            [pred      (lambda (x) `(eq ,val ,x))]

            [case*
             (lambda
               [() '(raise "no cases could be matched")]
               [clauses
                 (let [clause (first clauses)]
                   (assert-args
                     (is-vector clause)
                       (str "case clause must be a vector: " clause)
                     (= 2 (length clause))
                       (str "case clause must be paired: " clause))
                   (let ([test   (clause 0)]
                         [branch (clause 1)]
                         [next   (rest clauses)])
                     `(if ,(if (is-list test)
                               (pred-list test)
                               (pred test))
                           ,branch
                           ,(apply case* next))))])])
    `(let [,val ,expr]
        ,(apply case* cases))))

(define-macro if-let
  [(binding then) `(if-let ,binding ,then '())]
  [(binding then else)
    (assert-args
      (is-vector binding)
        (str "binding vector must be supplied: " binding)
      (= 2 (length binding))
        (str "binding vector must be paired: " binding))
    (let ([sym  (binding 0)]
          [test (binding 1)])
      `(let [,sym ,test]
            (if ,sym ,then ,else)))])

(define-macro when-let
  [(binding form)   `(if-let ,binding ,form)]
  [(binding . body) `(if-let ,binding (begin ,@body))])

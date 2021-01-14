;;;; ale core: branching

(define not is-false)

(define-macro unless
  [(test)           nil]
  [(test then)      `(if ,test nil ,then)]
  [(test then else) `(if ,test ,else ,then)])

(define-macro when
  [(test)         nil]
  [(test form)    `(if ,test ,form nil)]
  [(test . forms) `(if ,test (begin ,@forms) nil)])

(define-macro when-not
  [(test)         nil]
  [(test form)    `(if ,test nil ,form)]
  [(test . forms) `(if ,test nil (begin ,@forms))])

(define-macro cond
  [() nil]
  [clauses
   (let [clause (first clauses)]
     (assert-args
       (and (is-vector clause)
            (= 2 (length clause)))
       (str "invalid cond clause: " clause))
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
  (let-rec ([val       (gensym "val")                    ]
            [pred-list (lambda (l) `(or ,@(map! pred l)))]
            [pred      (lambda (x) `(eq ,val ,x))        ]

            [case*
             (lambda
               [() '(raise "no cases could be matched")]
               [clauses
                  (let [clause (first clauses)]
                    (assert-args
                      (and (is-vector clause)
                           (= 2 (length clause)))
                      (str "invalid case clause: " clause))
                    (let ([test   (clause 0)    ]
                          [branch (clause 1)    ]
                          [next   (rest clauses)])
                      `(if ,(if (is-list test)
                                (pred-list test)
                                (pred test))
                           ,branch
                           ,(apply case* next))))])])
    `(let [,val ,expr]
       ,(apply case* cases))))

(define-macro if-let
  [(binding then)
     `(if-let ,binding ,then nil)]
  [(binding then else)
     (assert-args
       (and (is-vector binding)
            (= 2 (length binding)))
       (str "invalid if-let binding: " binding))
     (let ([sym  (binding 0)]
           [test (binding 1)])
       `(let [,sym ,test]
           (if ,sym ,then ,else)))])

(define-macro when-let
  [(binding form)
     `(if-let ,binding ,form)]
  [(binding . body)
     `(if-let ,binding (begin ,@body))])

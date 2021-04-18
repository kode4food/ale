;;;; ale core: threading

(define (thread-seq->list target)
  (unless (list? target)
          (list target)
          target))

(define-macro ->
  [(value) value]
  [(value . forms)
     (let* ([l (thread-seq->list (first forms))]
            [f (first l)                       ]
            [r (rest l)                        ])
       `(-> (,f ,value ,@r) ,@(rest forms)))])

(define-macro ->>
  [(value) value]
  [(value . forms)
     (let* ([l (thread-seq->list (first forms))]
            [f (first l)                       ]
            [r (rest l)                        ])
       `(->> (,f ,@r ,value) ,@(rest forms)))])

(define-macro some->
  [(value) value]
  [(value . forms)
     (let* ([l (thread-seq->list (first forms))]
            [f (first l)                       ]
            [r (rest l)                        ])
       `(let [val# ,value]
          (when-not (null? val#)
                    (some-> (,f val# ,@r) ,@(rest forms)))))])

(define-macro some->>
  [(value) value]
  [(value . forms)
     (let* ([l (thread-seq->list (first forms))]
            [f (first l)                       ]
            [r (rest l)                        ])
       `(let [val# ,value]
          (when-not (null? val#)
                    (some->> (,f ,@r val#) ,@(rest forms)))))])

(define-macro (let-> binding . forms)
  (assert-args
    (is-binding-clause binding) "binding clause must be a paired vector"
    (!empty? forms)             "at least one threaded form is required")
  (let ([step  (thread-seq->list (first forms))]
        [next  (rest forms)                    ]
        [name  (binding 0)                     ]
        [value (binding 1)                     ])
    (if (empty? next)
        `(let [,name ,value] ,step)
        `(let [,name ,value]
           (let-> [,name ,step] ,@next)))))

(define (make-cond-clause threader)
  (lambda (clause)
    (assert-args
      (vector? clause)      "clause must be a vector"
      (= 2 (length clause)) "clause must be paired")
    (let ([pred (clause 0)]
          [form (clause 1)])
      `((lambda (val) (if ,pred (,threader val ,form) val))))))

(define-macro cond->
  [(value) value]
  [(value . clauses)
     `(-> ,value
          ,@(map (make-cond-clause '->) clauses))])

(define-macro cond->>
  [(value) value]
  [(value . clauses)
     `(-> ,value
          ,@(map (make-cond-clause '->>) clauses))])
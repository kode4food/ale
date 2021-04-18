;;;; ale core: exceptions

(letfn [(lambda-rec is-call (sym clause)
           (and (local? sym)
                (list? clause)
                (eq sym (first clause))))

        (lambda-rec is-catch-binding (form)
           (and (vector? form)
                (= 2 (length form))
                (local? (form 1))))

        (lambda-rec is-catch (clause parsed)
           (and (is-call 'catch clause)
                (is-catch-binding (nth clause 1))
                (!seq? (:block parsed))))

        (lambda-rec is-finally (clause parsed)
           (and (is-call 'finally clause)
                (!seq? (:catch parsed))
                (!seq? (:block parsed))))

        (lambda-rec is-expr (clause parsed)
           (!or (is-call 'catch clause)
                (is-call 'finally clause)))

        (lambda-rec try-append (parsed keyword clause)
           (conj parsed [keyword (conj (keyword parsed) clause)]))

        (lambda-rec try-prepend (parsed keyword clause)
           (conj parsed [keyword (cons clause (keyword parsed))]))

        (lambda-rec try-parse (clauses)
           (unless (seq? clauses)
                   {:block '() :catch '() :finally []}
                   (let* ([f (first clauses)]
                          [r (rest clauses) ]
                          [p (try-parse r)  ])
                     (cond
                       [(is-catch f p)   (try-prepend p :catch f)             ]
                       [(is-finally f p) (try-append p :finally f)            ]
                       [(is-expr f p)    (try-prepend p :block f)             ]
                       [:else            (raise "malformed try-catch-finally")]))))

        (lambda-rec try-catch-predicate (pred err-sym)
           (let* ([l (thread-seq->list pred)]
                  [f (first l)              ]
                  [r (rest l)               ])
             (cons f (cons err-sym r))))

        (lambda-rec try-catch-branch (clauses err-sym)
           (assert-args
             (seq? clauses) "catch branch not paired")
           (lazy-seq
             (let* ([clause (first clauses)     ]
                    [var    ((clause 1) 0)      ]
                    [expr   (rest (rest clause))])
               (cons (list 'ale/let
                           [var err-sym]
                           [false (cons 'ale/begin expr)])
                     (try-catch-clauses (rest clauses) err-sym)))))

        (lambda-rec try-catch-clauses (clauses err-sym)
           (lazy-seq
             (when (seq clauses)
               (let* ([clause (first clauses)]
                      [pred   ((clause 1) 1) ])
                 [(try-catch-predicate pred err-sym)
                  (try-catch-branch clauses err-sym)]))))

        (lambda-rec try-body (clauses)
           `(lambda () [false (begin ,@clauses)]))

        (lambda-rec try-catch (clauses)
           (let [err (gensym "err")]
             `(lambda (,err)
                (cond
                  ,@(apply list (try-catch-clauses clauses err))
                  [:else [true ,err]]))))

        (lambda-rec try-catch-finally (parsed)
           (let ([block   (:block parsed)  ]
                 [recover (:catch parsed)  ]
                 [cleanup (:finally parsed)])
             (cond
               [(seq? cleanup)
                (let ([first# (rest (first cleanup))                 ]
                      [rest#  (conj parsed [:finally (rest cleanup)])])
                  `(defer
                     (lambda () ,(try-catch-finally rest#))
                     (lambda () ,@first#)))]

               [(seq? recover)
                `(let ([rec# (recover ,(try-body block) ,(try-catch recover))]
                       [err# (rec# 0)                                        ]
                       [res# (rec# 1)                                        ])
                   (if err# (raise res#) res#))]

               [(seq? block) `(begin ,@block)]

               [:else        nil])))]

  (define-macro (try . clauses)
    (try-catch-finally (try-parse clauses))))

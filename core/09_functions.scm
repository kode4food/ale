;;;; ale core: functions

(define-macro (letfn bindings . body)
  ((lambda-rec parse-bindings (out in)
      (if (seq in)
        (let* ([fnList (first in)           ]
               [fnSym  (first fnList)       ]
               [fnName (first (rest fnList))]
               [fnRest (rest (rest fnList)) ])
          (assert-args
            (and (list? fnList)
                 (or (eq fnSym 'lambda-rec)
                     (eq fnSym 'ale/lambda-rec))
                 (local? fnName))
            "bindings must contain named functions")
          (parse-bindings (append out [fnName fnList]) (rest in)))
        `(let-rec ,(seq->list out) ,@body)))
   [] bindings))

(define-lambda partial
  [(func) func]
  [(func . first-args)
     (assert-args
       (apply? func) "partial requires a function")
     (lambda rest-args
       (apply func (apply append* (cons first-args rest-args))))])

(define-macro comp
  [() identity]
  [(func) func]
  [(func . funcs)
     (let* ([args        (gensym "args")        ]
            [inner       (list 'apply func args)]
            [first-outer (first funcs)          ]
            [rest-outer  (rest funcs)           ])
       (letfn [(lambda-rec outer (func args rest-funcs)
                  (if (seq rest-funcs)
                      (outer (first rest-funcs)
                             (list func args)
                             (rest rest-funcs))
                      (list func args)))]
         `(lambda ,args
            ,(outer first-outer inner rest-outer))))])

(define-macro (juxt . funcs)
  (let [args (gensym "args")]
    `(lambda ,args
       [,@(map (lambda (f) (list 'apply f args)) funcs)])))

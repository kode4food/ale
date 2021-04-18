;;;; ale core: binding

(define (is-binding-clause clause)
  (and (is-vector clause)
       (= 2 (length clause))
       (is-local (clause 0))))

(define (make-bindings value)
  (let-rec [is-bindings
            (lambda (value)
              (or (is-binding-clause value)
                  (and (is-list value)
                       (or (is-empty value)
                           (and (is-binding-clause (first value))
                                (is-bindings (rest value)))))))]
    (assert-args
      (is-bindings value) (str "invalid binding: " value))
    (if (is-vector value)
        (list value)
        value)))

(define-macro (let* bindings . body)
  (let [b (make-bindings bindings)]
    (let ([binding (first b)]
          [next    (rest b) ])
      (let ([name  (binding 0)]
            [value (binding 1)])
        (if (is-empty next)
            `(let [,name ,value] ,@body)
            `(let [,name ,value]
               (let* ,next ,@body)))))))

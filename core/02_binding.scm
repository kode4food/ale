;;;; ale core: binding

(define-macro (let* bindings . body)
  (assert-args
    (is-vector bindings)        "let* bindings must be a vector"
    (is-even (length bindings)) "let* bindings must be paired")
  (let [name  (bindings 0)
        value (bindings 1)]
    (if (> (length bindings) 2)
        `(let [,name ,value]
           (let* ,(rest (rest bindings)) ,@body))
        `(let [,name ,value] ,@body))))

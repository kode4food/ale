;;;; ale core: binding

(define-macro (let* bindings . body)
  (assert-args
    (is-list bindings) "let* bindings must be a list")
  (let ([binding (first bindings)]
        [next    (rest bindings)])
    (let ([name  (binding 0)]
          [value (binding 1)])
      (if (is-empty next)
          `(let [,name ,value] ,@body)
          `(let [,name ,value]
             (let* ,next ,@body))))))

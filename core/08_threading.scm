;;;; ale core: threading

(define (thread-to-list target)
  (unless (list? target)
          (list target)
          target))

(define-macro ->
  (lambda
    ([value] value)
    ([value . forms]
      (let* [l (thread-to-list (first forms))
            f (first l)
            r (rest l)]
        `(-> (,f ,value ,@r) ,@(rest forms))))))

(define-macro ->>
  (lambda
    ([value] value)
    ([value . forms]
      (let* [l (thread-to-list (first forms))
            f (first l)
            r (rest l)]
        `(->> (,f ,@r ,value) ,@(rest forms))))))

(define-macro some->
  (lambda
    ([value] value)
    ([value . forms]
      (let* [l (thread-to-list (first forms))
            f (first l)
            r (rest l)]
        `(let [val# ,value]
          (when-not (null? val#)
                    (some-> (,f val# ,@r) ,@(rest forms))))))))

(define-macro some->>
  (lambda
    ([value] value)
    ([value . forms]
      (let* [l (thread-to-list (first forms))
            f (first l)
            r (rest l)]
        `(let [val# ,value]
          (when-not (null? val#)
                    (some->> (,f ,@r val#) ,@(rest forms))))))))

(define-macro as->
  (lambda
    ([value name] value)
    ([value name . forms]
      (let [l (thread-to-list (first forms))]
        `(let [,name ,value]
          (as-> ,l ,name ,@(rest forms)))))))

(define (make-cond-clause sym)
  (lambda [clause]
    (let [pred (nth clause 0)
          form (nth clause 1)]
      `((lambda [val] (if ,pred (,sym val ,form) val))))))

(define-macro cond->
  (lambda
    ([value] value)
    ([value . clauses]
      (assert-args
        (even? (length clauses)) "clauses must be paired")
      `(-> ,value
          ,@(map (make-cond-clause ->) (partition 2 clauses))))))

(define-macro cond->>
  (lambda
    ([value] value)
    ([value . clauses]
      (assert-args
        (even? (length clauses)) "clauses must be paired")
      `(-> ,value
          ,@(map (make-cond-clause ->>) (partition 2 clauses))))))

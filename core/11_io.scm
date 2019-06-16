;;;; ale core: i/o

(declare *in* *out* *err*)

(define *space*   "\s")
(define *newline* "\n")

(define (pr-map-with-null func seq)
  (map (lambda [val] (if (null? val) val (func val)))
       seq))

(define (pr . forms)
  (let [mapped (pr-map-with-null str! forms)]
    (when (seq mapped)
          (. *out* :write (first mapped)))
    (when (seq mapped)
          (for-each [elem (rest mapped)]
                    (. *out* :write *space* elem)))))

(define (prn . forms)
  (apply pr forms)
  (. *out* :write *newline*))

(define (print . forms)
  (let [mapped (pr-map-with-null str forms)]
    (when (seq mapped)
          (. *out* :write (first mapped)))
    (when (seq mapped)
          (for-each [elem (rest mapped)]
                    (. *out* :write *space* elem)))))

(define (println . forms)
  (apply print forms)
  (. *out* :write *newline*))

(define (paired-vector? val)
  (and (vector? val) (pair? val)))

(define (with-open-close val)
  (let [c (:close val)]
    (if (apply? c) c no-op)))

(define-macro (with-open bindings . body)
  (assert-args
    (paired-vector? bindings) "with-open bindings must be a key-value vector")
  (cond
    (= (length bindings) 0)
    `(do ,@body)

    (>= (length bindings) 2)
    `(let [,(bindings 0) ,(bindings 1)
           close#        (with-open-close ,(bindings 0))]
       (try
         (with-open [,@(rest (rest bindings))] ,@body)
         (finally (close#))))))

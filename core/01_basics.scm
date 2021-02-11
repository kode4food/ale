;;;; ale core: basics

(define* *pos-inf* (/ 1.0 0.0))
(define* *neg-inf* (/ -1.0 0.0))

(define* true #t)
(define* false #f)
(define* nil '())

;; syntax-quoting requires it
(define* concat!
  (lambda colls
    (let-rec [concat-inner
              (lambda (colls head)
                (if (is-empty colls)
                    (apply list head)
                    (let ([f (first colls)]
                          [r (rest colls) ])
                      (if (is-empty f)
                          (concat-inner r head)
                          (concat-inner (cons (rest f) r)
                                        (append head (first f)))))))]
      (concat-inner colls []))))

(define* label
  (macro
    (lambda (name form)
      `(let-rec [,name ,form] ,name))))

(define* is-cons-or-list
  (lambda (value)
    (if (is-cons value)
        true
        (is-list value))))

(let [make-macro
      (lambda (quoter)
        (macro
          (lambda forms
            (let ([f (car forms)]
                  [r (cdr forms)])
              (if (is-cons-or-list f)
                  (quoter (car f) (cons (cdr f) r))
                  (quoter f r))))))]

  (define* define-lambda
    (make-macro
      (lambda (name body)
        `(define* ,name
           (label ,name (lambda ,@body))))))

  (define* define-macro
    (make-macro
      (lambda (name body)
        `(define* ,name
           (label ,name (macro (lambda ,@body))))))))

(define-macro assert-args
  [() nil]
  [(clause)
     (raise "assert-args clauses must be paired")]
  [clauses
     `(if ,(clauses 0)
          (assert-args ,@(rest (rest clauses)))
          (raise ,(clauses 1)))])

(define-macro (lambda-rec name . forms)
  (if (is-local name)
      `(label ,name (lambda ,@forms))
      `(lambda ,name ,@forms)))

(define-macro (!eq value . comps)
  `(not (eq ,value ,@comps)))

(define-macro and
  [() true]
  [(clause) clause]
  [clauses
     `(let [and# ,(clauses 0)]
        (if and#
            (and ,@(rest clauses))
            and#))])

(define-macro (!and . clauses)
  `(not (and ,@clauses)))

(define-macro or
  [() nil]
  [(clause) clause]
  [clauses
     `(let [or# ,(clauses 0)]
        (if or#
            or#
            (or ,@(rest clauses))))])

(define-macro (!or . clauses)
  `(not (or ,@clauses)))

(define-macro (declare . names)
  (let-rec [declare
            (lambda (names)
              (if (is-empty names)
                  '()
                  (cons (list 'ale/declare* (first names))
                        (declare (rest names)))))]
    `(begin ,@(declare names))))

(define-macro (define . body)
  (let [value (first body)]
    (assert-args
      (or (is-local value)
          (is-cons-or-list value))
      (str "invalid define: " body))
    (if (is-local value)
        `(define* ,@body)
        `(define-lambda ,@body))))

(define (is-even value)
  (= (mod value 2) 0))

(define (is-odd value)
  (= (mod value 2) 1))

(define (is-true value)
  (if value true false))

(define (is-false value)
  (if value false true))

(define (inc value)
  (+ value 1))

(define (dec value)
  (- value 1))

(define (no-op . _))

(define (identity value) value)

(define (constantly value)
  (lambda _ value))

(define-macro (: target method . args)
  `((get ,target ,method) ,@args))

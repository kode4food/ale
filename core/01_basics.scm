;;;; ale core: basics

(declare *env* *args*)

(def *pos-inf* (/ 1.0 0.0))
(def *neg-inf* (/ -1.0 0.0))

;; syntax-quoting requires it
(def concat!
  (lambda colls
    (let-rec [concat-inner
              (lambda (colls head)
                (if (is-empty colls)
                    (apply list head)
                    (let ([f (first colls)]
                          [r (rest colls)])
                      (if (is-empty f)
                          (concat-inner r head)
                          (concat-inner (cons (rest f) r)
                                        (append head (first f)))))))]
      (concat-inner colls []))))

(def label
  (macro
    (lambda (name form)
      `(let-rec [,name ,form] ,name))))

(let [parse-define
      (lambda (forms)
        (let ([f (car forms)]
              [r (cdr forms)])
          (if (is-pair f)
              [(car f) (cons (cdr f) r)]
              [f r])))]

  (def define-lambda
    (label define-lambda
      (macro
        (lambda forms
          (let [parsed (parse-define forms)]
            (let ([name (parsed 0)]
                  [body (parsed 1)])
              `(def ,name
                 (label ,name (lambda ,@body)))))))))

  (def define-macro
    (label define-macro
      (macro
        (lambda forms
          (let [parsed (parse-define forms)]
            (let ([name (parsed 0)]
                  [body (parsed 1)])
              `(def ,name
                 (label ,name (macro (lambda ,@body)))))))))))

(define-macro assert-args
  [() '()]
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

(define-macro (define . body)
  (if (is-local (car body))
      `(def ,@body)
      `(define-lambda ,@body)))

(define-macro (!eq value . comps)
  `(not (eq ,value ,@comps)))

(define-macro and
  [() #t]
  [(clause) clause]
  [clauses
    `(let [and# ,(clauses 0)]
      (if and#
          (and ,@(rest clauses))
          and#))])

(define-macro (!and . clauses)
  `(not (and ,@clauses)))

(define-macro or
  [() '()]
  [(clause) clause]
  [clauses
    `(let [or# ,(clauses 0)]
      (if or#
          or#
          (or ,@(rest clauses))))])

(define-macro (!or . clauses)
  `(not (or ,@clauses)))

(define (is-even value)
  (= (mod value 2) 0))

(define (is-odd value)
  (= (mod value 2) 1))

(define (is-true value)
  (if value #t #f))

(define (is-false value)
  (if value #f #t))

(define (inc value)
  (+ value 1))

(define (dec value)
  (- value 1))

(define (no-op . _))

(define (identity value) value)

(define (constantly value)
  (lambda _ value))

(define-macro (^ target method . args)
  `((get ,target ,method) ,@args))

;;;; ale core: predicates

(let-rec ([pred-apply
           (lambda (func args)
             (if (is-empty args)
                 true
                 (unless (func (first args)) false
                         (pred-apply func (rest args)))))]

          [define-pos
           (lambda (func name)
             (let [func-name (sym (str name "?"))]
               `(define ,func-name (lambda (f# . r#)
                  (,pred-apply ,func (cons f# r#))))))]

          [define-neg
           (lambda (func name)
             (let [func-name (sym (str "!" name "?"))]
               `(define ,func-name (lambda (f# . r#)
                  (,pred-apply (lambda (value) (not (,func value)))
                               (cons f# r#))))))])

  (define-macro (define-predicate func name)
    `(begin ,(define-pos func name)
            ,(define-neg func name))))

(define (is-null value)
  (eq value nil))

(define (is-zero value)
  (= value 0))

(define-predicate is-appender "append")
(define-predicate is-apply "apply")
(define-predicate is-atom "atom")
(define-predicate is-boolean "boolean")
(define-predicate is-cons "cons")
(define-predicate is-counted "counted")
(define-predicate is-empty "empty")
(define-predicate is-even "even")
(define-predicate is-false "false")
(define-predicate is-indexed "indexed")
(define-predicate is-keyword "keyword")
(define-predicate is-list "list")
(define-predicate is-local "local")
(define-predicate is-macro "macro")
(define-predicate is-mapped "mapped")
(define-predicate is-nan "nan")
(define-predicate is-neg-inf "-inf")
(define-predicate is-null "null")
(define-predicate is-number "number")
(define-predicate is-object "object")
(define-predicate is-odd "odd")
(define-predicate is-pair "pair")
(define-predicate is-pos-inf "inf")
(define-predicate is-promise "promise")
(define-predicate is-qualified "qualified")
(define-predicate is-resolved "resolved")
(define-predicate is-reversible "reversible")
(define-predicate is-seq "seq")
(define-predicate is-special "special")
(define-predicate is-string "string")
(define-predicate is-symbol "symbol")
(define-predicate is-true "true")
(define-predicate is-vector "vector")
(define-predicate is-zero "zero")

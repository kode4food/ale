;;;; ale core: predicates

(letrec [pred-apply
         (lambda [func args]
           (if (is-empty args)
               #t
               (unless (func (first args)) #f
                       (pred-apply func (rest args)))))

         def-predicate-pos
         (lambda [func name]
           (let [func-name (sym (str name "?"))]
             `(define (,func-name ,@'(first . rest))
                (,pred-apply ,func (cons ,'first ,'rest)))))

         def-predicate-neg
         (lambda [func name]
           (let [func-name (sym (str "!" name "?"))]
             `(define (,func-name ,@'(first . rest))
                (,pred-apply (lambda [value#] (not (,func value#)))
                             (cons ,'first ,'rest)))))]

  (defmacro def-predicate [func name]
    `(do ,(def-predicate-pos func name)
         ,(def-predicate-neg func name))))

(def-predicate is-appender "append")
(def-predicate is-apply "apply")
(def-predicate is-object "object")
(def-predicate is-atom "atom")
(def-predicate is-boolean "boolean")
(def-predicate is-counted "counted")
(def-predicate is-delivered "delivered")
(def-predicate is-deque "deque")
(def-predicate is-empty "empty")
(def-predicate is-even "even")
(def-predicate is-indexed "indexed")
(def-predicate is-keyword "keyword")
(def-predicate is-list "list")
(def-predicate is-local "local")
(def-predicate is-macro "macro")
(def-predicate is-mapped "mapped")
(def-predicate is-nan "nan")
(def-predicate is-neg-inf "-inf")
(def-predicate is-null "null")
(def-predicate is-number "number")
(def-predicate is-odd "odd")
(def-predicate is-pair "pair")
(def-predicate is-pos-inf "inf")
(def-predicate is-promise "promise")
(def-predicate is-qualified "qualified")
(def-predicate is-reversible "reversible")
(def-predicate is-seq "seq")
(def-predicate is-special "special")
(def-predicate is-string "string")
(def-predicate is-symbol "symbol")
(def-predicate is-vector "vector")
(def-predicate is-true "true")
(def-predicate is-false "false")

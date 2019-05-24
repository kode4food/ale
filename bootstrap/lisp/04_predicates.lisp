;;;; ale bootstrap: predicates

(defn pred-apply
  [func args]
  (if (is-empty args) true
      (unless (func (first args)) false
              (pred-apply func (rest args)))))

(defmacro def-predicate-pos
  [func name]
  (let [func-name (sym (str name "?"))]
    `(defn ~func-name [~'first & ~'rest]
       (pred-apply ~func (cons ~'first ~'rest)))))

(defmacro def-predicate-neg
  [func name]
  (let [func-name (sym (str "!" name "?"))]
    `(defn ~func-name [~'first & ~'rest]
       (not (pred-apply ~func (cons ~'first ~'rest))))))

(defmacro def-predicate
  [func name]
  `(do
     (def-predicate-pos ~func ~name)
     (def-predicate-neg ~func ~name)))

(def-predicate is-appender "append")
(def-predicate is-apply "apply")
(def-predicate is-assoc "assoc")
(def-predicate is-atom "atom")
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
(def-predicate is-nil "nil")
(def-predicate is-odd "odd")
(def-predicate is-paired "paired")
(def-predicate is-pos-inf "inf")
(def-predicate is-promise "promise")
(def-predicate is-qualified "qualified")
(def-predicate is-reversible "reversible")
(def-predicate is-seq "seq")
(def-predicate is-special "special")
(def-predicate is-str "str")
(def-predicate is-symbol "symbol")
(def-predicate is-vector "vector")
(def-predicate is-true "true")
(def-predicate is-false "false")


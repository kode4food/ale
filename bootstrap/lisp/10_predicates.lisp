;;;; ale bootstrap: predicates

(defn predicate-lazy-seq
  [func coll]
  (lazy-seq
   (if (seq coll)
     (let [f (func (first coll))]
       (if f
         (cons f (predicate-lazy-seq func (rest coll)))
         '(false)))
     '(true))))

(defmacro def-predicate-pos
  [func name]
  (let [func-name (sym (str name "?"))]
    `(defn ~func-name [~'first & ~'rest]
       (last
        (predicate-lazy-seq ~func (cons ~'first ~'rest))))))

(defmacro def-predicate-neg
  [func name]
  (let [func-name (sym (str "!" name "?"))]
    `(defn ~func-name [~'first & ~'rest]
       (let [nf# (fn [x] (not (~func x)))]
         (last (predicate-lazy-seq nf# (cons ~'first ~'rest)))))))

(defmacro def-predicate
  [func name]
  `(do
     (def-predicate-pos ~func ~name)
     (def-predicate-neg ~func ~name)))

(def-predicate is-atom "atom")
(def-predicate is-nil "nil")
(def-predicate is-str "str")
(def-predicate is-seq "seq")
(def-predicate is-empty "empty")
(def-predicate is-len "len")
(def-predicate is-indexed "indexed")
(def-predicate is-assoc "assoc")
(def-predicate is-mapped "mapped")
(def-predicate is-list "list")
(def-predicate is-vector "vector")

(def-predicate is-promise "promise")

(def-predicate is-pos-inf "inf")
(def-predicate is-neg-inf "-inf")
(def-predicate is-nan "nan")

(def-predicate is-macro "macro")
(def-predicate is-apply "apply")
(def-predicate is-special "special")

(def-predicate is-keyword "keyword")
(def-predicate is-symbol "symbol")
(def-predicate is-local "local")
(def-predicate is-qualified "qualified")

(def-predicate is-even "even")
(def-predicate is-odd "odd")
(def-predicate is-paired "paired")

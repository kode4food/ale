;;;; ale bootstrap: basics

(declare *env* *args*)
(declare syntax-quote)

(def *pos-inf* (/ 1.0 0.0))
(def *neg-inf* (/ -1.0 0.0))

(defmacro defn
  [name & forms]
  `(def ~name (fn ~name ~@forms)))

(defmacro error
  [& clauses]
  `(raise (assoc ~@clauses)))

(defmacro panic
  [& clauses]
  `(raise (error ~@clauses)))

(defmacro eq
  [value & comps]
  `(is-eq ~value ~@comps))

(defmacro !eq
  [value & comps]
  `(not (is-eq ~value ~@comps)))

(defn is-even
  [value]
  (= (mod value 2) 0))

(defn is-odd
  [value]
  (= (mod value 2) 1))

(defn is-paired
  [value]
  (is-even (len value)))

(defn inc
  [value]
  (+ value 1))

(defn dec
  [value]
  (- value 1))

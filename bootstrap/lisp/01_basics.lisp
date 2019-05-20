;;;; ale bootstrap: basics

(declare *env* *args*)

(def *pos-inf* (/ 1.0 0.0))
(def *neg-inf* (/ -1.0 0.0))

(defmacro defn
  [name & forms]
  `(def ~name (fn ~name ~@forms)))

(defn error
  [& clauses]
  (raise (apply assoc clauses)))

(defn panic
  [& clauses]
  (raise (apply error clauses)))

(defmacro !eq
  [value & comps]
  `(not (eq ~value ~@comps)))

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

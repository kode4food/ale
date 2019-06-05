;;;; ale core: basics

(declare *env* *args*)

(def *pos-inf* (/ 1.0 0.0))
(def *neg-inf* (/ -1.0 0.0))

;; syntax-quoting requires it
(def concat!
  (lambda [& colls]
    (letrec [concat'
             (lambda [colls head]
               (if (is-empty colls)
                   (apply list head)
                   (let [f (first colls)
                         r (rest colls)]
                     (if (is-empty f)
                         (concat' r head)
                         (concat' (cons (rest f) r)
                                  (append head (first f)))))))]
      (concat' colls []))))

(def defmacro
  (letrec [defmacro
           (macro [name & forms]
             `(def ~name
                (letrec [~name (macro ~@forms)]
                  ~name)))]
    defmacro))

(defmacro fn
  [name & forms]
  (if (is-local name)
    `(letrec [~name (lambda ~@forms)] ~name)
    `(lambda ~name ~@forms)))

(defmacro defn
  [name & forms]
  `(def ~name (fn ~name ~@forms)))

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

(defn is-true
  [value]
  (if value true false))

(defn is-false
  [value]
  (if value false true))

(defn inc
  [value]
  (+ value 1))

(defn dec
  [value]
  (- value 1))

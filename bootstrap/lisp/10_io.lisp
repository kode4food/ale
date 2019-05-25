;;;; ale bootstrap: i/o

(declare *in* *out* *err*)

(def *space*   "\s")
(def *newline* "\n")

(defn pr-map-with-nil
  [func seq]
  (map
   (fn [val] (if (nil? val) val (func val)))
   seq))

(defn pr [& forms]
  (let [mapped (pr-map-with-nil str! forms)]
    (if (seq mapped)
      (. *out* :write (first mapped)))
    (for-each [elem (rest mapped)]
              (. *out* :write *space* elem))))

(defn prn [& forms]
  (apply pr forms)
  (. *out* :write *newline*))

(defn print [& forms]
  (let [mapped (pr-map-with-nil str forms)]
    (if (seq mapped)
      (. *out* :write (first mapped)))
    (for-each [elem (rest mapped)]
              (. *out* :write *space* elem))))

(defn println [& forms]
  (apply print forms)
  (. *out* :write *newline*))

(defn paired-vector?
  [val]
  (and (vector? val) (paired? val)))

(defn with-open-close
  [val]
  (let [c (:close val)]
    (if (apply? c) c no-op)))

(defmacro with-open [bindings & body]
  (assert-args
   (paired-vector? bindings) "with-open bindings must be a key-value vector")
  (cond
    (= (len bindings) 0)
    `(do ~@body)
    (>= (len bindings) 2)
    `(let [~(bindings 0) ~(bindings 1)
           close#        (with-open-close ~(bindings 0))]
       (try
         (with-open [~@(rest (rest bindings))] ~@body)
         (finally (close#))))))

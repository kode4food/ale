;;;; ale bootstrap: i/o

(declare *in* *out* *err*)

(def *space*   "\s")
(def *newline* "\n")

(defn pr-map-with-nil
  [func seq]
  (map
   (fn [val] (if (is-nil val) val (func val)))
   seq))

(defn pr [& forms]
  (let [seq (pr-map-with-nil str! forms)]
    (if (is-seq seq)
      (. *out* :write (first seq)))
    (for-each [elem (rest seq)]
              (. *out* :write *space* elem))))

(defn prn [& forms]
  (apply pr forms)
  (. *out* :write *newline*))

(defn print [& forms]
  (let [seq (pr-map-with-nil str forms)]
    (if (is-seq seq)
      (. *out* :write (first seq)))
    (for-each [elem (rest seq)]
              (. *out* :write *space* elem))))

(defn println [& forms]
  (apply print forms)
  (. *out* :write *newline*))

(defn paired-vector?
  [val]
  (and (is-vector val) (= (mod (len val) 2) 0)))

(defn with-open-close
  [val]
  (let [c (:close val)]
    (if (is-apply c) c no-op)))

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

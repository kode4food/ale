;;;; ale bootstrap: lazy sequences

(defmacro lazy-seq
  [& body]
  `(lazy-seq* (fn [] ~@body)))

(defmacro to-assoc
  [& seqs]
  `(apply assoc (concat ~@seqs)))

(defmacro to-list
  [& seqs]
  `(apply list (concat ~@seqs)))

(defmacro to-vector
  [& seqs]
  `(apply vector (concat ~@seqs)))

(defn take-
  [count coll]
  (lazy-seq
   (when-let [s (and (> count 0) (seq coll))]
     (cons (first s)
           (take- (dec count) (rest s))))))

(defn take
  [count coll]
  (assert-args
   (= (mod count 1) 0) "count must be an integer"
   (>= count 0)        "count must be non-negative"
   (seq coll)          "coll must be a sequence")
  (take- count coll))

(defn take-while
  [pred coll]
  (lazy-seq
   (when-let [s (seq coll)]
     (when (pred (first s))
       (cons (first s)
             (take-while pred (rest s)))))))

(defmacro for-each
  [seq-exprs & body]
  (assert-args
   (is-vector seq-exprs) "for-each bindings must be a vector"
   (is-paired seq-exprs) "for-each bindings must be paired")
  (let [name# (seq-exprs 0)
        seq#  (seq-exprs 1)]
    (if (> (len seq-exprs) 2)
      (let [rest# (rest (rest seq-exprs))]
        `(for-each* ~seq# (fn [~name#] (for-each ~rest# ~@body))))
      `(for-each* ~seq# (fn [~name#] ~@body)))))

(defmacro for
  [seq-exprs & body]
  `(generate
    (for-each ~seq-exprs (emit (do ~@body)))))

(defn partition
  ([count coll] (partition count count coll))
  ([count step coll]
   (lazy-seq
    (when (is-seq coll)
      (cons (to-list (take count coll))
            (partition count step (drop step coll)))))))

(defn range
  ([]
   (range 0 nil 1))
  ([last]
   (range 0 last (if (> last 0) 1 -1)))
  ([first last]
   (if (> last first)
     (range first last 1)
     (range last first -1)))
  ([first last step]
   (let [cmp (cond (is-nil last) (constantly true)
                   (< step 0)    >
                   :else         <)]
     (if (cmp first last)
       (cons first (lazy-seq (range (+ first step) last step)))
       []))))

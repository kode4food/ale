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

(defn len
  [coll]
  (assert-args
   (seq? coll) "coll must be a sequence")
  ((fn len'
     [coll prev]
     (if (sized? coll)
       (+ prev (size coll))
       (if (seq coll)
         (len' (rest coll) (inc prev))
         prev)))
   coll 0))

(defn last
  [coll]
  (assert-args
   (seq? coll) "coll must be a sequence")
  ((fn last'
     [coll prev]
     (if (and (sized? coll) (indexed? coll))
       (let [s (size coll)]
         (if (> s 0)
           (nth coll (dec (size coll)))
           prev))
       (if (seq coll)
         (let [f (first coll)
               r (rest coll)]
           (last' r f))
         prev)))
   coll nil))

(defn take
  [count coll]
  (assert-args
   (= (mod count 1) 0) "count must be an integer"
   (seq? coll)         "coll must be a sequence")
  ((fn take'
     [count coll]
     (lazy-seq
      (when-let [s (and (> count 0) (seq coll))]
        (cons (first s)
              (take' (dec count) (rest s))))))
   count coll))

(defn take-while
  [pred coll]
  (lazy-seq
   (when-let [s (seq coll)]
     (let [fs (first s)]
       (when (pred fs)
         (cons fs (take-while pred (rest s))))))))

(defn drop
  [count coll]
  (assert-args
   (= (mod count 1) 0) "count must be an integer"
   (seq? coll)         "coll must be a sequence")
  (lazy-seq
   ((fn drop'
      [count coll]
      (let [s (seq coll)]
        (if (and s (> count 0))
          (drop' (dec count) (rest s))
          s)))
    count coll)))

(defmacro for-each
  [seq-exprs & body]
  (assert-args
   (vector? seq-exprs) "for-each bindings must be a vector"
   (paired? seq-exprs) "for-each bindings must be paired")
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
    (when (seq? coll)
      (cons (to-list (take count coll))
            (partition count step (drop step coll)))))))

(defn range
  ([]     (range 0 nil 1))
  ([last] (range 0 last (if (> last 0) 1 -1)))
  ([first last]
   (if (> last first)
     (range first last 1)
     (range last first -1)))
  ([first last step]
   (let [cmp (cond (nil? last) (constantly true)
                   (< step 0)  >
                   :else       <)]
     (if (cmp first last)
       (cons first (lazy-seq (range (+ first step) last step)))
       []))))

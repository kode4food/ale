;;;; ale bootstrap: lazy sequences

(defmacro lazy-seq
  [& body]
  `(lazy-seq* (fn [] ~@body)))

(defn to-assoc
  [& seqs]
  (apply assoc (apply concat seqs)))

(defn to-list
  [& seqs]
  (apply list (apply concat seqs)))

(defn to-vector
  [& seqs]
  (apply vector (apply concat seqs)))

(defn append* [coll & values]
  ((fn append' [coll values]
     (if (seq values)
       (append' (append coll (first values)) (rest values))
       coll))
   coll values))

(defn prepend* [coll & values]
  ((fn prepend' [coll values]
     (if (seq values)
       (prepend' (cons (first values) coll) (rest values))
       coll))
   coll values))

(defn conj [coll & values]
  (if (append? coll)
    (apply append* (cons coll values))
    (apply prepend* (cons coll values))))

(defn len!
  [coll]
  ((fn len'
     [coll prev]
     (if (counted? coll)
       (+ prev (len coll))
       (if (seq coll)
         (len' (rest coll) (inc prev))
         prev)))
   coll 0))

(defn last
  [coll]
  (let [s (len coll)]
    (when (> s 0)
      (nth coll (dec (len coll))))))

(defn last!
  [coll]
  ((fn last'
     [coll prev]
     (if (and (counted? coll) (indexed? coll))
       (let [s (len coll)]
         (if (> s 0)
           (nth coll (dec (len coll)))
           prev))
       (if (seq coll)
         (let [f (first coll)
               r (rest coll)]
           (last' r f))
         prev)))
   coll nil))

(defn reverse!
  [coll]
  (if (is-reversible coll)
    (reverse coll)
    ((fn reverse' [coll target]
       (if (seq coll)
         (reverse' (rest coll) (cons (first coll) target))
         target))
     coll ())))

(defn take
  [count coll]
  ((fn take'
     [count coll]
     (lazy-seq
      (if (and (> count 0) (!empty? coll))
        (cons (first coll) (take' (dec count) (rest coll)))
        ())))
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
  (lazy-seq
   ((fn drop'
      [count coll]
      (if (> count 0)
        (drop' (dec count) (rest coll))
        coll))
    count coll)))
(defn for-each' [func coll]
  (when (seq coll)
    (do (func (first coll))
        (for-each' func (rest coll)))))
(defmacro for-each
  [seq-exprs & body]
  (assert-args
   (vector? seq-exprs) "for-each bindings must be a vector"
   (paired? seq-exprs) "for-each bindings must be paired")
  (let [name# (seq-exprs 0)
        seq#  (seq-exprs 1)]
    (if (> (len seq-exprs) 2)
      (let [rest# (rest (rest seq-exprs))]
        `(for-each' (fn [~name#] (for-each ~rest# ~@body)) ~seq#))
      `(for-each' (fn [~name#] ~@body) ~seq#))))

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

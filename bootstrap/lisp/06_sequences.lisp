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

(defn cp-rotate-row [row orig-row]
  (if (seq row)
    (let [res (rest row)]
      (if (seq res)
        [false res]
        [true orig-row]))
    [true orig-row]))

(defn cp-rotate-rest [rest orig]
  (let [f (first rest) fo (first orig)
        r (rest rest)  ro (rest orig)]
    (if (seq r)
      (let [res (cp-rotate-rest r ro)]
        (if (res 0)
          (let [rr (cp-rotate-row f fo)]
            [(rr 0) (cons (rr 1) (res 1))])
          [false (cons f (res 1))]))
      (let [res (cp-rotate-row f fo)]
        [(res 0) (list (res 1))]))))

(defn cp-rotate [work orig]
  (let [res (cp-rotate-rest work orig)]
    (unless (res 0) (res 1))))

(defn cartesian-product
  [& seqs]
  ((fn iter [work]
     (lazy-seq
      (let [f (to-vector (map first work))
            r (cp-rotate work seqs)]
        (if r
          (cons f (iter r))
          (list f)))))
   seqs))

(defn for-seq-exprs
  ([name seq]
   [[name] [seq]])
  ([name seq & rest]
   (let [res (apply for-seq-exprs rest)]
     [(cons name (res 0))
      (cons seq (res 1))])))

(defmacro for
  [seq-exprs & body]
  (assert-args
   (vector? seq-exprs) "for-each bindings must be a vector"
   (paired? seq-exprs) "for-each bindings must be paired")
  (let [split (apply for-seq-exprs seq-exprs)
        names# (split 0)
        seqs#  (split 1)]
    `(map
      (fn [args] (apply (fn ~names# ~@body) args))
      (cartesian-product ~@seqs#))))

(defmacro for-each
  [seq-exprs & body]
  `(last! (for ~seq-exprs ~@body)))

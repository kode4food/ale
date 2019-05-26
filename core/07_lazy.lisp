;;;; ale core: lazy sequences

(defmacro lazy-seq
  [& body]
  `(lazy-seq* (fn [] ~@body)))

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
  ([count coll]
   (partition count count coll))

  ([count step coll]
   (lazy-seq
    (when (seq? coll)
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
   (let [cmp (cond (nil? last) (constantly true)
                   (< step 0)  >
                   :else       <)]
     (if (cmp first last)
       (cons first (lazy-seq (range (+ first step) last step)))
       []))))

(defn map
  ([func coll]
   ((fn map' [coll]
      (lazy-seq
       (when (seq coll)
         (cons (func (first coll))
               (map' (rest coll))))))
    coll))

  ([func coll & colls]
   ((fn parallel' [colls]
      (lazy-seq
       (when (apply true? (map !empty? colls))
         (let [f (to-vector (map first colls))
               r (map rest colls)]
           (cons (apply func f) (parallel' r))))))
    (cons coll colls))))

(defn filter
  [func coll]
  (lazy-seq
   ((fn filter' [coll]
      (when (seq coll)
        (let [f (first coll) r (rest coll)]
          (if (func f)
            (cons f (filter func r))
            (filter' r)))))
    coll)))

(defn cartesian-product
  [& colls]
  (let [rotate-row
        (fn rotate-row [row orig-row]
          (if (seq row)
            (let [res (rest row)]
              (if (seq res)
                [false res]
                [true orig-row]))
            [true orig-row]))

        rotate-rest
        (fn rotate-rest [rest orig]
          (let [f (first rest) fo (first orig)
                r (rest rest)  ro (rest orig)]
            (if (seq r)
              (let [res (rotate-rest r ro)]
                (if (res 0)
                  (let [rr (rotate-row f fo)]
                    [(rr 0) (cons (rr 1) (res 1))])
                  [false (cons f (res 1))]))
              (let [res (rotate-row f fo)]
                [(res 0) (list (res 1))]))))

        rotate
        (fn rotate [work orig]
          (let [res (rotate-rest work orig)]
            (unless (res 0) (res 1))))]
    ((fn iter [work]
       (lazy-seq
        (let [f (to-vector (map first work))
              r (rotate work colls)]
          (if r
            (cons f (iter r))
            (list f)))))
     colls)))

(defmacro for
  [seq-exprs & body]
  (assert-args
   (vector? seq-exprs) "for-each bindings must be a vector"
   (paired? seq-exprs) "for-each bindings must be paired")
  (let [args (gensym "args")

        split-bindings
        (fn split-bindings
          ([idx name coll]
           [(list name (list args idx))
            (list coll)])
          ([idx name coll & rest]
           (let [res (apply split-bindings (cons (inc idx) rest))]
             [(cons* (res 0) (list args idx) name)
              (cons coll (res 1))])))

        split (apply split-bindings (cons 0 seq-exprs))
        bind# (to-vector (split 0))
        seqs# (split 1)]
    `(map
      (fn [~args] (let ~bind# ~@body))
      (cartesian-product ~@seqs#))))

(defmacro for-each
  [seq-exprs & body]
  `(last! (for ~seq-exprs ~@body)))

(defn concat [& colls]
  ((fn concat' [colls]
     (lazy-seq
      (when (seq colls)
        (let [f (first colls)
              r (rest colls)]
          (if (seq f)
            (cons (first f)
                  (concat' (cons (rest f) r)))
            (concat' r))))))
   colls))

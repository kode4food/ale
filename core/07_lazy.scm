;;;; ale core: lazy sequences

(define-macro (lazy-seq . body)
  `(lazy-seq* (lambda [] ,@body)))

(define (take count coll)
  ((fn take-inner [count coll]
     (lazy-seq
       (if (and (> count 0) (!empty? coll))
           (cons (first coll) (take-inner (dec count) (rest coll)))
           '())))
    count coll))

(define (take-while pred coll)
  (lazy-seq
    (when-let [s (seq coll)]
      (let [fs (first s)]
        (when (pred fs)
              (cons fs (take-while pred (rest s))))))))

(define (drop count coll)
  (lazy-seq
    ((fn drop-inner [count coll]
       (if (> count 0)
           (drop-inner (dec count) (rest coll))
           coll))
      count coll)))

(defn partition
  ([count coll]
    (partition count count coll))

  ([count step coll]
    (lazy-seq
      (when (seq coll)
            (cons (to-list (take count coll))
                  (partition count step (drop step coll)))))))

(defn range
  ([]
    (range 0 '() 1))

  ([last]
    (range 0 last (if (> last 0) 1 -1)))

  ([first last]
    (if (> last first)
        (range first last 1)
        (range last first -1)))

  ([first last step]
    (let [cmp (cond (null? last) (constantly #t)
                    (< step 0)  >
                    :else       <)]
      (if (cmp first last)
          (cons first (lazy-seq (range (+ first step) last step)))
          []))))

(defn map
  ([func coll]
    ((fn map-single [coll]
       (lazy-seq
         (when (seq coll)
               (cons (func (first coll))
                     (map-single (rest coll))))))
      coll))

  ([func coll . colls]
    ((fn map-parallel [colls]
       (lazy-seq
         (when (apply true? (map !empty? colls))
               (let [f (to-vector (map first colls))
                     r (map rest colls)]
                 (cons (apply func f) (map-parallel r))))))
      (cons coll colls))))

(define (filter func coll)
  (lazy-seq
    ((fn filter-inner [coll]
       (when (seq coll)
             (let [f (first coll)
                   r (rest coll)]
               (if (func f)
                   (cons f (filter func r))
                   (filter-inner r)))))
      coll)))

(define-macro (for-each seq-exprs . body)
  `(last! (for ,seq-exprs ,@body)))

(define (concat . colls)
  ((fn concat-inner [colls]
     (lazy-seq
       (when (seq colls)
             (let [f (first colls)
                   r (rest colls)]
               (if (seq f)
                   (cons (first f)
                         (concat-inner (cons (rest f) r)))
                   (concat-inner r))))))
     colls))

(define (zip . colls)
  (apply map list colls))

(define (mapcat func . colls)
  (apply concat (apply map func colls)))

(define-macro (for seq-exprs . body)
  (assert-args
    (vector? seq-exprs)        "for-each bindings must be a vector"
    (even? (length seq-exprs)) "for-each bindings must be paired"
    (!empty? seq-exprs)        "at least one binding pair is required")
  (let [sym  (seq-exprs 0)
        expr (seq-exprs 1)
        next (rest (rest seq-exprs))]
    (if (= (length seq-exprs) 2)
        `(map (lambda [,sym] ,@body) (seq! ,expr))
        `(mapcat (lambda [,sym] (for ,next ,@body)) (seq! ,expr)))))

(define-macro (cartesian-product . colls)
  (let* [sym-gen  (lambda [x] (gensym (str "cp" x)))
         let-syms (take (length colls) (map sym-gen (range)))
         for-syms (take (length colls) (map sym-gen (range)))
         let-vals (zip let-syms colls)
         for-vals (zip for-syms let-syms)
         let-bind (to-vector (apply concat let-vals))
         for-bind (to-vector (apply concat for-vals))]
    `(let ,let-bind (for ,for-bind (vector ,@for-syms)))))

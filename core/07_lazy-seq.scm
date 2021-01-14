;;;; ale core: lazy sequences

(define-macro (lazy-seq . body)
  `(lazy-seq* (lambda () ,@body)))

(define (take count coll)
  ((lambda-rec take-inner (count coll)
      (lazy-seq
        (if (and (> count 0)
                 (!empty? coll))
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
    ((lambda-rec drop-inner (count coll)
        (if (> count 0)
            (drop-inner (dec count) (rest coll))
            coll))
     count coll)))

(define-lambda partition
  [(count coll)
     (partition count count coll)]

  [(count step coll)
     (lazy-seq
       (when (seq coll)
             (cons (seq->list (take count coll))
                   (partition count step (drop step coll)))))])

(define-lambda range
  [()
     (range 0 nil 1)]

  [(last)
     (range 0 last (if (> last 0) 1 -1))]

  [(first last)
     (if (> last first)
         (range first last 1)
         (range last first -1))]

  [(first last step)
     (let [cmp
           (cond [(null? last) (constantly true)]
                 [(< step 0)   >                ]
                 [:else        <                ])]
       (if (cmp first last)
           (cons first (lazy-seq (range (+ first step) last step)))
           []))])

(define-lambda map
  [(func coll)
     ((lambda-rec map-single (coll)
         (lazy-seq
           (when (seq coll)
                 (cons (func (first coll))
                       (map-single (rest coll))))))
      coll)]

  [(func coll . colls)
     ((lambda-rec map-parallel (colls)
         (lazy-seq
           (when (apply true? (map !empty? colls))
                 (let ([f (seq->vector (map first colls))]
                       [r (map rest colls)               ])
                   (cons (apply func f) (map-parallel r))))))
      (cons coll colls))])

(define (filter func coll)
  (lazy-seq
    ((lambda-rec filter-inner (coll)
        (when (seq coll)
              (let ([f (first coll)]
                    [r (rest coll) ])
                (if (func f)
                    (cons f (filter func r))
                    (filter-inner r)))))
     coll)))

(define-macro (for-each seq-exprs . body)
  `(last! (for ,seq-exprs ,@body)))

(define (concat . colls)
  ((lambda-rec concat-inner (colls)
      (lazy-seq
        (when (seq colls)
              (let ([f (first colls)]
                    [r (rest colls) ])
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
  (let [b (make-bindings seq-exprs)]
    (assert-args
      (!empty? b) "at least one for binding is required")
    (let* ([this (first b)]
           [next (rest b) ]
           [sym  (this 0) ]
           [expr (this 1) ])
      (if (empty? next)
          `(map (lambda (,sym) ,@body) (seq! ,expr))
          `(mapcat (lambda (,sym) (for ,next ,@body)) (seq! ,expr))))))

(define-macro (cartesian-product . colls)
  (let* ([sym-gen  (lambda (x) (gensym (str "cp" x)))         ]
         [let-syms (take (length colls) (map sym-gen (range)))]
         [let-vals (map seq->vector (zip let-syms colls))     ]
         [let-bind (seq->list let-vals)                       ]
         [for-syms (take (length colls) (map sym-gen (range)))]
         [for-vals (map seq->vector (zip for-syms let-syms))  ]
         [for-bind (seq->list for-vals)                       ])
    `(let ,let-bind (for ,for-bind (list ,@for-syms)))))

;;;; ale core: standard sequences

(define (seq value)
  (if (and (is-seq value) (!empty? value))
      value
      '()))

(define (seq! value)
  (if (is-seq value)
      (if (!empty? value) value '())
      (raise (str "can't treat " value " as a sequence"))))

(define (to-object . colls)
  (apply object (apply concat! colls)))

(define (to-list . colls)
  (apply list (apply concat! colls)))

(define (to-vector . colls)
  (apply vector (apply concat! colls)))

(define (append* coll . values)
  ((fn append-inner [coll values]
     (if (seq values)
         (append-inner (append coll (first values)) (rest values))
         coll))
    coll values))

(define (cons* coll . values)
  ((fn cons-inner [coll values]
     (if (seq values)
         (cons-inner (cons (first values) coll) (rest values))
         coll))
    coll values))

(define (conj coll . values)
  (if (append? coll)
      (apply append* (cons coll values))
      (apply cons* (cons coll values))))

(define (length! coll)
  ((fn length-inner [coll prev]
     (if (counted? coll)
         (+ prev (length coll))
         (if (seq coll)
             (length-inner (rest coll) (inc prev))
             prev)))
    coll 0))

(define nth!
  (let [scan
        (fn scan [coll pos handle]
          (if (seq coll)
              (if (> pos 0)
                  (scan (rest coll) (dec pos) handle)
                  (first coll))
              (handle)))]
    (fn nth!
      ([coll pos]
        (if (indexed? coll)
            (nth coll pos)
            (scan coll pos (lambda [] (raise "index out of bounds")))))
      ([coll pos default]
        (if (indexed? coll)
            (nth coll pos default)
            (scan coll pos (lambda [] default)))))))

(define (last coll)
  (let [s (length coll)]
    (when (> s 0)
          (nth coll (dec (length coll))))))

(define (last! coll)
  ((fn last-inner [coll prev]
     (if (and (counted? coll) (indexed? coll))
         (let [s (length coll)]
           (if (> s 0)
               (nth coll (dec (length coll)))
               prev))
         (if (seq coll)
             (let [f (first coll)
                   r (rest coll)]
               (last-inner r f))
             prev)))
    coll '()))

(defn reduce
  ([func init coll]
    ((fn reduce-inner [init coll]
       (if (seq coll)
           (reduce-inner (func init (first coll)) (rest coll))
           init))
      init coll))

  ([func coll]
    (if (seq coll)
        (reduce func (first coll) (rest coll))
        (func))))

(define (reverse! coll)
  (if (is-reversible coll)
      (reverse coll)
      (reduce conj '() coll)))

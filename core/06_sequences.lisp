;;;; ale core: standard sequences

(define (to-assoc & colls)
  (apply assoc (apply concat! colls)))

(define (to-list & colls)
  (apply list (apply concat! colls)))

(define (to-vector & colls)
  (apply vector (apply concat! colls)))

(define (append* coll & values)
  ((fn append' [coll values]
     (if (seq values)
         (append' (append coll (first values)) (rest values))
         coll))
    coll values))

(define (cons* coll & values)
  ((fn cons' [coll values]
     (if (seq values)
         (cons' (cons (first values) coll) (rest values))
         coll))
    coll values))

(define (conj coll & values)
  (if (append? coll)
      (apply append* (cons coll values))
      (apply cons* (cons coll values))))

(define (len! coll)
  ((fn len'
     [coll prev]
     (if (counted? coll)
         (+ prev (len coll))
         (if (seq coll)
             (len' (rest coll) (inc prev))
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
  (let [s (len coll)]
    (when (> s 0)
          (nth coll (dec (len coll))))))

(define (last! coll)
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

(define (reverse! coll)
  (if (is-reversible coll)
      (reverse coll)
      ((fn reverse' [coll target]
         (if (seq coll)
             (reverse' (rest coll) (cons (first coll) target))
             target))
        coll ())))

(defn reduce
  ([func init coll]
    ((fn reduce' [init coll]
       (if (seq coll)
           (reduce' (func init (first coll)) (rest coll))
           init))
      init coll))

  ([func coll]
    (if (seq coll)
        (reduce func (first coll) (rest coll))
        (func))))

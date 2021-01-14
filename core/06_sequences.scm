;;;; ale core: standard sequences

(define (seq value)
  (if (and (seq? value)
           (!empty? value))
      value
      nil))

(define (seq! value)
  (if (seq? value)
      (if (!empty? value) value nil)
      (raise (str "value can't act as a sequence: " value))))

(define (seq->object . colls)
  (apply object (apply concat! colls)))

(define (seq->list . colls)
  (apply list (apply concat! colls)))

(define (seq->vector . colls)
  (apply vector (apply concat! colls)))

(define (append* coll . values)
  ((lambda-rec append-inner (coll values)
      (if (seq values)
          (append-inner (append coll (first values)) (rest values))
          coll))
   coll values))

(define (cons* coll . values)
  ((lambda-rec cons-inner (coll values)
      (if (seq values)
          (cons-inner (cons (first values) coll) (rest values))
          coll))
   coll values))

(define (conj coll . values)
  (if (append? coll)
      (apply append* (cons coll values))
      (apply cons* (cons coll values))))

(define (length! coll)
  ((lambda-rec length-inner (coll prev)
      (if (counted? coll)
          (+ prev (length coll))
          (if (seq coll)
              (length-inner (rest coll) (inc prev))
              prev)))
   coll 0))

(define nth!
  (let-rec [scan
            (lambda (coll pos missing)
              (if (seq coll)
                  (if (> pos 0)
                      (scan (rest coll) (dec pos) missing)
                      (first coll))
                  (missing)))]
    (lambda-rec nth!
       [(coll pos)
          (if (indexed? coll)
              (nth coll pos)
              (scan coll pos (lambda () (raise "index out of bounds"))))]
       [(coll pos default)
          (if (indexed? coll)
              (nth coll pos default)
              (scan coll pos (lambda () default)))])))

(define (last coll)
  (let [s (length coll)]
    (when (> s 0)
          (nth coll (dec (length coll))))))

(define (last! coll)
  ((lambda-rec last-inner (coll prev)
      (if (and (counted? coll)
               (indexed? coll))
          (let [s (length coll)]
            (if (> s 0)
                (nth coll (dec (length coll)))
                prev))
          (if (seq coll)
              (let ([f (first coll)]
                    [r (rest coll) ])
                (last-inner r f))
              prev)))
   coll '()))

(define-lambda reduce
  [(func init coll)
     ((lambda-rec reduce-inner (init coll)
         (if (seq coll)
             (reduce-inner (func init (first coll)) (rest coll))
             init))
      init coll)]

  [(func coll)
     (if (seq coll)
         (reduce func (first coll) (rest coll))
         (func))])

(define (reverse! coll)
  (if (reversible? coll)
      (reverse coll)
      (reduce conj '() coll)))

;; ale core - alist

(define (assoc key coll)
  (when (seq coll)
    (let [elem (first coll)]
      (unless (and (seq elem) (eq key (first elem)))
              (assoc key (rest coll))
              elem))))

(define (alist->object coll)
  (apply object
         (mapcat (lambda (elem) [(first elem) (rest elem)])
                 (filter seq coll))))

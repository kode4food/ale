;;;; ale core: os

(declare *env* *args*)

(define-macro (time . forms)
  `(let* ([start#  (current-time) ]
          [result# (begin ,@forms)]
          [end#    (current-time) ]
          [dur#    (- end# start#)])
     (if (< dur# 1000)
         (println dur# "ns")
         (println (/ dur# 1000000.0) "ms"))
     result#))

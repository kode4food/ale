;;;; ale core: os

(defmacro time
  [& forms]
  `(let* [start#  (current-time)
          result# (do ,@forms)
          end#    (current-time)
          dur#    (- end# start#)]
     (if (< dur# 1000)
         (println dur# "ns")
         (println (/ dur# 1000000) "ms"))
     result#))

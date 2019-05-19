;;;; ale bootstrap: concurrency

(defmacro go
  [& body]
  `(go* (fn [] ~@body)))

(defmacro generate
  [& body]
  `(let [chan#  (chan)
         close# (:close chan#)
         emit   (:emit chan#)]
     (go
       (let [result# (do ~@body)]
         (close#)
         result#))
     (:seq chan#)))

(defmacro future
  [& body]
  `(let [promise# (promise)]
     (go (promise# (do ~@body)))
     promise#))

;;;; ale core: concurrency

(defmacro go
  [& body]
  `(go* (lambda [] ,@body)))

(defmacro generate
  [& body]
  `(let* [chan#  (chan)
          close# (:close chan#)
          emit   (:emit chan#)]
     (go
       (let [result# (do ,@body)]
         (close#)
         result#))
     (:seq chan#)))

(defmacro future
  [& body]
  `(let [promise# (promise)]
     (go (promise# (do ,@body)))
     promise#))

(defn spawn
  ([func]
    (spawn func 16))
  ([func mbox-size]
    (spawn func mbox-size no-op))
  ([func mbox-size monitor]
    (let* [channel (chan mbox-size)
           mailbox (:seq channel)
           sender  (:emit channel)]
      (go
        (recover (lambda [] (func mailbox))
                 monitor))
                 sender)))

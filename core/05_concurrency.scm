;;;; ale core: concurrency

(define-macro (go . body)
  `(go* (lambda [] ,@body)))

(define-macro (generate . body)
  `(let* [chan#  (chan)
          close# (:close chan#)
          emit   (:emit chan#)]
     (go
       (let [result# (begin ,@body)]
         (close#)
         result#))
     (:seq chan#)))

(define-macro (future . body)
  `(let [promise# (promise)]
     (go (promise# (begin ,@body)))
     promise#))

(define-macro (delay expr)
  `(lazy-seq [,expr]))

(define-macro (force value)
  `(first ,value))

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

;;;; ale core: concurrency

(define-macro (go . body)
  `(go* (lambda () ,@body)))

(define-macro (future . body)
  `(let [promise# (promise (lambda () ,@body))]
     (go (promise#))
     promise#))

(define-macro (delay value)
  `(promise (lambda () ,value)))

(define (force value)
  (if (is-promise value)
      (value)
      (raise "attempt to force a non-promise")))

(define-macro (generate . body)
  `(let* ([chan#  (chan)]
          [close# (:close chan#)]
          [emit   (:emit chan#)])
     (go
       (let [result# (begin ,@body)]
         (close#)
         result#))
     (:seq chan#)))

(define-lambda spawn
  [(func)
     (spawn func 16)]
  [(func mbox-size)
     (spawn func mbox-size no-op)]
  [(func mbox-size monitor)
     (let* ([channel (chan mbox-size)]
            [mailbox (:seq channel)  ]
            [sender  (:emit channel) ])
       (go
         (recover (lambda () (func mailbox))
                   monitor))
                   sender)])

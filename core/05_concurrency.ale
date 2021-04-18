;;;; ale core: concurrency

(define-macro (go . body)
  `(go* (lambda () ,@body)))

(define-macro (future . body)
  `(let [promise# (promise (lambda () ,@body))]
     (go (promise#))
     promise#))

(define-macro (delay . body)
  `(promise (lambda () ,@body)))

(define (force value)
  (if (promise? value)
      (value)
      value))

(define-macro (lazy . body)
  `(delay
    (let [body-result# (begin ,@body)]
      ((lambda-rec resolve (result)
         (if (promise? result)
             (resolve (result))
             result))
       body-result#))))

(define-macro (generate . body)
  `(let* ([chan#  (chan)        ]
          [close# (:close chan#)]
          [emit   (:emit chan#) ])
     (go
       (let [result# (begin ,@body)]
         (close#)
         result#))
     (:seq chan#)))

;; spawn an actor. The provided func accepts a single mailbox argument
;; that is a channel sequence. Returns a sender function that can send
;; messages to the mailbox. Default mailbox size before send operations
;; begin to block is 16 messages
(define-lambda spawn
  [(func)
     (spawn func 16)]
  [(func mbox-size)
     (spawn func mbox-size no-op)]
  [(func mbox-size monitor)
     (let* ([channel (chan mbox-size)]
            [mailbox (:seq channel)  ]
            [sender  (:emit channel) ])
       (go (recover (lambda () (func mailbox))
                    monitor))
       sender)])

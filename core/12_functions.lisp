;;;; ale core: functions

(defmacro letfn
  [bindings & body]
  ((fn parse-bindings [out in]
     (if (seq in)
       (let* [fnList (first in)
              fnSym  (first fnList)
              fnName (first (rest fnList))
              fnRest (rest (rest fnList))]
         (assert-args
           (and (is-list fnList)
                (or (eq fnSym 'fn) (eq fnSym 'ale/fn))
                (is-local fnName))
           "bindings must contain named functions")
         (parse-bindings (append (append out fnName) fnList) (rest in)))
       `(letrec ~out ~@body)))
   [] bindings))

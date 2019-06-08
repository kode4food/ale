;;;; ale core: functions

(defmacro letfn
  [bindings . body]
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
       `(letrec ,out ,@body)))
   [] bindings))

(defn partial
  ([func] func)
  ([func . first-args]
    (assert-args
      (is-apply func) "partial requires a function")
    (lambda rest-args
      (apply func (apply append* (cons first-args rest-args))))))

(defmacro comp
  ([] identity)
  ([func] func)
  ([func . funcs]
    (let* [args        (gensym "args")
           inner       (list 'apply func args)
           first-outer (first funcs)
           rest-outer  (rest funcs)]
      (letfn [(fn outer [func args rest-funcs]
                (if (seq rest-funcs)
                    (outer (first rest-funcs)
                           (list func args)
                           (rest rest-funcs))
                    (list func args)))]
        `(lambda ,args
           ,(outer first-outer inner rest-outer))))))

(defmacro juxt funcs
  (let [args (gensym "args")]
    `(lambda ,args
       [,@(map (lambda [f] (list 'apply f args)) funcs)])))

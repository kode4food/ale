;;;; ale core: calling

(defn partial
  ([func] func)
  ([func & first-args]
   (assert-args
    (is-apply func) "partial requires a function")
   (fn [& rest-args]
     (apply func (apply append* (cons first-args rest-args))))))

(defmacro comp
  ([] identity)
  ([func] func)
  ([func & funcs]
   (let [args        (gensym "args")
         inner       (list 'apply func args)
         first-outer (first funcs)
         rest-outer  (rest funcs)

         outer
         (fn outer
           [func args rest-funcs]
           (if (seq rest-funcs)
             (outer (first rest-funcs) (list func args) (rest rest-funcs))
             (list func args)))]
     `(fn [& ~args]
        ~(outer first-outer inner rest-outer)))))

(defmacro juxt
  [& funcs]
  (let [args (gensym "args")]
    `(fn [& ~args]
       [~@(map (fn [f] (list 'apply f args)) funcs)])))

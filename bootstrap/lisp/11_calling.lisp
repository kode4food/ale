;;;; ale bootstrap: calling

(defn partial
  ([func] func)
  ([func & first-args]
   (assert-args
    (is-apply func) "partial requires a function")
   (fn [& rest-args]
     (apply func (apply append* (cons first-args rest-args))))))

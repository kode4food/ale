;;;; ale core: threading

(defn thread-to-list
  [target]
  (unless (list? target)
          (list target)
          target))

(defmacro ->
  ([value] value)
  ([value & forms]
   (let [l (thread-to-list (first forms))
         f (first l)
         r (rest l)]
     `(-> (~f ~value ~@r) ~@(rest forms)))))

(defmacro ->>
  ([value] value)
  ([value & forms]
   (let [l (thread-to-list (first forms))
         f (first l)
         r (rest l)]
     `(->> (~f ~@r ~value) ~@(rest forms)))))

(defmacro some->
  ([value] value)
  ([value & forms]
   (let [l (thread-to-list (first forms))
         f (first l)
         r (rest l)]
     `(let [val# ~value]
        (when-not (nil? val#)
          (some-> (~f val# ~@r) ~@(rest forms)))))))

(defmacro some->>
  ([value] value)
  ([value & forms]
   (let [l (thread-to-list (first forms))
         f (first l)
         r (rest l)]
     `(let [val# ~value]
        (when-not (nil? val#)
          (some->> (~f ~@r val#) ~@(rest forms)))))))

(defmacro as->
  ([value name] value)
  ([value name & forms]
   (let [l (thread-to-list (first forms))]
     `(let [~name ~value]
        (as-> ~l ~name ~@(rest forms))))))

(defn make-cond-clause
  [sym]
  (fn [clause]
    (let [pred (nth clause 0)
          form (nth clause 1)]
      `((fn [val] (if ~pred (~sym val ~form) val))))))

(defmacro cond->
  ([value] value)
  ([value & clauses]
   (assert-args
    (even? (len clauses)) "clauses must be paired")
   `(-> ~value
        ~@(map (make-cond-clause ->) (partition 2 clauses)))))

(defmacro cond->>
  ([value] value)
  ([value & clauses]
   (assert-args
    (even? (len clauses)) "clauses must be paired")
   `(-> ~value
        ~@(map (make-cond-clause ->>) (partition 2 clauses)))))

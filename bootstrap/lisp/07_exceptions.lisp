;;;; ale bootstrap: exceptions

(defn is-call
  [sym clause]
  (and (is-local sym)
       (is-list clause)
       (eq sym (first clause))))

(defn is-catch-binding
  [form]
  (and (is-vector form)
       (= 2 (len form))
       (is-local (form 1))))

(defn is-catch
  [clause parsed]
  (and (is-call 'catch clause)
       (is-catch-binding (nth clause 1))
       (not (is-seq (:block parsed)))))

(defn is-finally
  [clause parsed]
  (and (is-call 'finally clause)
       (not (is-seq (:catch parsed)))
       (not (is-seq (:block parsed)))))

(defn is-expr
  [clause parsed]
  (!or (is-call 'catch clause)
       (is-call 'finally clause)))

(defn try-append
  [parsed keyword clause]
  (conj parsed [keyword (conj (keyword parsed) clause)]))

(defn try-prepend
  [parsed keyword clause]
  (conj parsed [keyword (cons clause (keyword parsed))]))

(defn try-parse
  [clauses]
  (unless (is-seq clauses)
          {:block () :catch () :finally []}
          (let [f (first clauses)
                r (rest clauses)
                p (try-parse r)]
            (cond
              (is-catch f p)   (try-prepend p :catch f)
              (is-finally f p) (try-append  p :finally f)
              (is-expr f p)    (try-prepend p :block f)
              :else            (raise "malformed try-catch-finally")))))

(defn try-catch-predicate
  [pred err-sym]
  (let [l (thread-to-list pred)
        f (first l)
        r (rest l)]
    (cons f (cons err-sym r))))

(declare try-catch-clauses)

(defn try-catch-branch
  [clauses err-sym]
  (assert-args (is-seq clauses) "catch branch not paired")
  (lazy-seq
   (let [clause (first clauses)
         var    ((clause 1) 0)
         expr   (rest (rest clause))]
     (cons
      (list 'ale/let
            [var err-sym]
            [false (cons 'ale/do expr)])
      (try-catch-clauses (rest clauses) err-sym)))))

(defn try-catch-clauses
  [clauses err-sym]
  (lazy-seq
   (when (is-seq clauses)
     (let [clause (first clauses)
           pred   ((clause 1) 1)]
       (cons
        (try-catch-predicate pred err-sym)
        (try-catch-branch clauses err-sym))))))

(defn try-body
  [clauses]
  `(fn [] [false (do ~@clauses)]))

(defn try-catch
  [clauses]
  (let [err (gensym "err")]
    `(fn [~err]
       (cond
         ~@(apply list (try-catch-clauses clauses err))
         :else [true ~err]))))

(defn try-catch-finally
  [parsed]
  (let [block   (:block parsed)
        recover (:catch parsed)
        cleanup (:finally parsed)]
    (cond (is-seq cleanup)
          (let [first# (rest (first cleanup))
                rest#  (conj parsed [:finally (rest cleanup)])]
            `(defer
               (fn [] ~(try-catch-finally rest#))
               (fn [] ~@first#)))

          (is-seq recover)
          `(let [rec# (recover ~(try-body block) ~(try-catch recover))
                 err# (rec# 0)
                 res# (rec# 1)]
             (if err# (raise res#) res#))

          (is-seq block) `(do ~@block)

          :else nil)))

(defmacro try
  [& clauses]
  (try-catch-finally (try-parse clauses)))
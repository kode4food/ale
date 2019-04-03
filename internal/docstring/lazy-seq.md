# (lazy-seq form**) produces a sequence that is evaluated lazily

## An Example:

  (defn fib-seq []
    (let [fib (fn fib[a b]
                (lazy-seq (cons a (fib b (+ a b)))))]
      (fib 0 1)))

  (for-each [x (take 300 (fib-seq))]
    (println x))

This example prints the first 300 fibonacci numbers.

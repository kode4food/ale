---
title: "lazy-seq"
description: "produces a sequence that is evaluated lazily"
names: ["lazy-seq"]
usage: "(lazy-seq form*)"
tags: ["sequence"]
---

#### An Example:

```scheme
(define (fib-seq)
  (let [fib (lambda-rec fib (a b)
              (lazy-seq (cons a (fib b (+ a b)))))]
    (fib 0 1)))

(for-each [x (take 300 (fib-seq))]
  (println x))
```

This example prints the first 300 fibonacci numbers.

---
title: "lazy-seq"
date: 2019-04-06T12:19:22+02:00
description: "produces a sequence that is evaluated lazily"
names: ["lazy-seq"]
usage: "(lazy-seq form*)"
tags: ["sequence"]
---

#### An Example:

```clojure
(defn fib-seq []
  (let [fib (fn fib[a b]
              (lazy-seq (cons a (fib b (+ a b)))))]
    (fib 0 1)))

(for-each [x (take 300 (fib-seq))]
  (println x))
```

This example prints the first 300 fibonacci numbers.

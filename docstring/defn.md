---
title: "defn"
date: 2019-04-06T12:19:22+02:00
description: "binds a namespace function"
names: ["defn", "define ()"]
usage: "(defn name (param*) form*) (define (name param*) form*)"
tags: ["function", "binding"]
---
Will bind a function by name to the current namespace.

#### An Example

~~~scheme
(define (fib i)
  (cond
    [(= i 0) 0]
    [(= i 1) 1]
    [(= i 2) 1]
    [:else   (+ (fib (- i 2)) (fib (- i 1)))]))
~~~

This example performs recursion with no tail call optimization, and no memoization. For a more performant and stack-friendly fibonacci sequence generation example, see the documentation of `lazy-seq`.

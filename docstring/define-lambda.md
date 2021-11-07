---
title: "define-lambda"
description: "binds a namespace function"
names: ["define-lambda"]
usage: "(define-lambda name (param*) form*) (define-lambda (name param*) form*)"
tags: ["function", "binding"]
---

Will bind a function by name to the current namespace.

#### An Example

```scheme
(define-lambda (fib i)
  (cond
    [(= i 0) 0]
    [(= i 1) 1]
    [(= i 2) 1]
    [:else   (+ (fib (- i 2)) (fib (- i 1)))]))
```

This example performs recursion with no tail call optimization, and no memoization. For a more performant and stack-friendly fibonacci sequence generation example, see the documentation of `lazy-seq`.

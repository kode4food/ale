---
title: "declare"
date: 2019-04-06T12:19:22+02:00
description: "forward declares a binding"
names: ["declare"]
usage: "(declare <name>+)"
tags: ["binding"]
---
Forward declares bindings. This means that the names will be known in the current namespace, but not yet assigned. This can be useful when two functions refer to one another.

#### An Example

~~~scheme
(declare is-odd-number)

(define (is-even-number n)
  (cond [(= n 0) #t]
        [:else   (is-odd-number (- n 1))]))

(define (is-odd-number n)
  (cond [(= n 0) #f]
        [:else   (is-even-number (- n 1))]))
~~~

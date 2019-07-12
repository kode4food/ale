---
title: "when"
date: 2019-04-06T12:19:22+02:00
description: "conditionally evaluates a block"
names: ["when", "when-not"]
usage: "(when pred form*)"
tags: ["conditional"]
---
If the evaluated predicate is truthy (not _#f_ (false) or the empty list), the forms are evaluated. Will evaluate each form in turn, returning the final evaluation as its result.

#### An Example

~~~scheme
(define x '(1 2 3 4 5 6 7 8))

(when (> (length x) 3)
  (println "x is big")
  (length x))
~~~

If the symbol `when-not` is used instead of `when`, then the predicate is evaluated and the block will be evaluated only if result is not truthy

---
title: "to-vector"
date: 2019-04-06T12:19:22+02:00
description: "converts sequences to a vector"
names: ["to-vector"]
usage: "(to-vector seq+)"
tags: ["sequence", "conversion"]
---
Will concatenate a set of sequences into a vector. Unlike the standard `concat` function, which is lazily computed, the result of this function will be materialized immediately.

#### An Example

~~~scheme
(define x
  (map (lambda (x) (* x 2))
  '(1 2 3 4)))
(to-vector '(1 2 3 4) x)
~~~
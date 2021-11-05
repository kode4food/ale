---
title: "seq->vector"
description: "converts sequences to a vector"
names: ["seq->vector"]
usage: "(seq->vector seq+)"
tags: ["sequence" "conversion"]
---

Will concatenate a set of sequences into a vector. Unlike the standard `concat` function, which is lazily computed, the result of this function will be materialized immediately.

#### An Example

```scheme
(define x
  (map (lambda (x) (* x 2))
  '(1 2 3 4)))
(seq->vector '(1 2 3 4) x)
```

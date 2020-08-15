---
title: "seq->list"
date: 2019-04-06T12:19:22+02:00
description: "converts sequences to a list"
names: ["seq->list"]
usage: "(seq->list seq+)"
tags: ["sequence", "conversion"]
---

Will concatenate a set of sequences into a list. Unlike the standard `concat` function, which is lazily computed, the result of this function will be materialized immediately.

#### An Example

```scheme
(define x [1 2 3 4])
(define y
  (map (lambda (x) (+ x 4))
  '(1 2 3 4)))
(seq->list x y)
```

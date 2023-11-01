---
title: "seq->object"
description: "converts sequences to an object"
names: ["seq->object"]
usage: "(seq->object seq+)"
tags: ["sequence", "conversion"]
---

Concatenate a set of sequences into an object (hash-map). Unlike the standard `concat` function, which is lazily computed, the result of this function will be materialized immediately.

#### An Example

```scheme
(define x [:name "ale" :age 0.3])
(define y '(:weight "light"))
(seq->object x y)
```

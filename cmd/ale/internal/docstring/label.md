---
title: "label"
description: "assigns a stable name to a procedure value"
names: ["label"]
usage: "(label name form)"
tags: ["function"]
---

Wraps a form with a stable name, primarily to support recursive procedure definitions and clearer generated values.

#### An Example

```scheme
(label fact
  (lambda (n)
    (if (zero? n) 1 (* n (fact (dec n))))))
```

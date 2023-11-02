---
title: "comp"
description: "composes a set of functions"
names: ["comp"]
usage: "(comp func*)"
tags: ["function"]
---

Returns a new function based on chained invocation of the provided functions, from left to right. The first composed function can accept multiple arguments, while any subsequent functions are applied with the result of the previous.

#### An Example

```scheme
(define mul2Add5 (comp (partial * 2) (partial + 5)))
(mul2Add5 10)
```

This example will return _25_ as though `(+ 5 (* 2 10))` were called.

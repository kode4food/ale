---
title: "partial"
description: "generates a function based on a partial apply"
names: ["partial"]
usage: "(partial func arg+)"
tags: ["function"]
---

Returns a new function whose initial arguments are pre-bound to those provided. When that function is invoked, any provided arguments will simply be appended to the pre-bound arguments before calling the original function.

#### An Example

```scheme
(define plus10 (partial + 4 6))
(plus10 9)
```

This example will return _19_ as though `(+ 4 6 9)` were called.

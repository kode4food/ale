---
title: "recover"
description: "evaluates a body and handles raised values"
names: ["recover"]
usage: "(recover body rescue)"
tags: ["exception"]
---

Invokes the zero-argument procedure `body`. If evaluation raises or panics, the runtime normalizes the error and passes the resulting Ale value to the single-argument procedure `rescue`. This is the procedural building block used by higher-level error macros such as `try`.

#### An Example

```scheme
(recover
  (thunk (raise "boom"))
  (lambda (err) (str "caught: " err)))
```

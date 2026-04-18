---
title: ":"
description: "invokes a looked-up method or value as a procedure"
names: [":"]
usage: "(: target method arg*)"
tags: ["sequence", "macro"]
---

Looks up `method` from `target` using `get`, then calls the resulting value as a procedure with any remaining arguments. It expands to a normal `get` followed by a call.

#### An Example

```scheme
(: {:inc (lambda (x) (+ x 1))} :inc 41)
```

This example returns `42`.

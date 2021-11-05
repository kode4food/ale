---
title: "apply"
usage: "(apply func <arg>* seq)"
names: ["apply"]
description: "applies a function to the provided arguments"
tags: ["function"]
---

Evaluates the provided sequence and applies the provided function to its values and any explicitly included arguments.

#### An Example

```scheme
(define x '(1 2 3))
(apply + x)
```

This example will return _6_.

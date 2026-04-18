---
title: "for"
description: "lazy sequence comprehensions"
names: ["for", "for-each"]
usage: "(for [binding seq]* body) (for-each [binding seq]* body)"
tags: ["sequence", "macro"]
---

`for` builds a lazy sequence by iterating over one or more bindings. `for-each` evaluates the same comprehension for side effects and returns the final result.

#### An Example

```scheme
(for ([x [1 2]]
      [y [10 20]])
  (+ x y))
```

This example lazily yields `11`, `21`, `12`, and `22`.

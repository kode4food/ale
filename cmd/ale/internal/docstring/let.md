---
title: "let"
description: "binds local values"
names: ["let", "let*", "let-rec"]
usage: "(let ([name expr]*) form*) (let* ([name expr]*) form*) (let-rec ([name expr]*) form*)"
tags: ["binding"]
---

Create a new local scope, evaluate the provided expressions, and then bind the resulting values to their respective names. The body is evaluated within that scope and returns the result of its final form.

The binding forms differ in how names can see one another. `let` performs bindings in parallel. `let*` performs bindings sequentially. `let-rec` allows recursive and mutually recursive references.

#### An Example

```scheme
(let ([x '(1 2 3 4)]
      [y [5 6 7 8] ])
  (concat x y))
```

This example will create a list called _x_ and a vector called _y_ and return the lazy concatenation of those sequences. Note that the two names do not exist outside the `let` form.

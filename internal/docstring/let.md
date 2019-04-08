---
title: "let"
date: 2019-04-06T12:19:22+02:00
description: "binds local values"
names: ["let"]
usage: "(let [<name form>+] form+)"
tags: ["binding"]
---
Will create a new local scope, binding the specified values to that scope by name. It will then evaluate the specified forms within that scope and return the result of the last evaluation.

#### An Example

```clojure
(let [x '(1 2 3 4)
      y [5 6 7 8]]
  (concat x y))
```

This example will create a list called *x* and a vector called *y* and return the lazy concatenation of those sequences. Note that the two names do not exist outside of the `let` form.

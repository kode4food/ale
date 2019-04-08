---
title: "if"
date: 2019-04-06T12:19:22+02:00
description: "performs simple branching"
names: ["if", "unless"]
usage: "(if pred then else?)"
tags: ["conditional"]
---
If the evaluated predicate is truthy (not _false_, not _nil_), the *then* form is evaluated and returned, otherwise the *else* form, if any, will be evaluated and returned.

#### An Example

```clojure
(def x '(1 2 3 4 5 6 7 8))

(if (> (len x) 3)
  "x is big"
  "x is small")
```

If the symbol `unless` is used instead of `if`, then the logical branching will be inverted.

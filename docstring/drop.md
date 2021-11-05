---
title: "drop"
description: "drops the first elements of a sequence"
names: ["drop"]
usage: "(drop count seq)"
tags: ["sequence" "comprehension"]
---

Will return a lazy sequence that excludes the first _count_ elements of the provided sequence. If the source sequence is shorter than the requested count, an empty list will be returned.

#### An Example

```scheme
(define x '(1 2 3 4))
(define y [5 6 7 8])
(drop 3 (concat x y))
```

This example will return the lazy sequence _(4 5 6 7 8)_.

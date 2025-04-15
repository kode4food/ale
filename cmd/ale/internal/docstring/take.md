---
title: "take"
description: "takes the first elements of a sequence"
names: ["take"]
usage: "(take count seq)"
tags: ["sequence", "comprehension"]
---

Return a lazy sequence of either _count_ or fewer elements from the beginning of the provided sequence. If the source sequence is shorter than the requested count, the resulting sequence will be truncated.

#### An Example

```scheme
(define x '(1 2 3 4))
(define y [5 6 7 8])
(take 6 (concat x y))
```

This example will return the lazy sequence _(1 2 3 4 5 6)_.

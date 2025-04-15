---
title: "last"
description: "returns the last element of the sequence"
names: ["last", "last!"]
usage: "(last seq) (last! seq)"
tags: ["sequence"]
---

This function will return the last element of the specified sequence, or the empty list if the sequence is empty. If the sequence is lazily computed, asynchronous, or otherwise incapable of returning a count, this function will raise an error.

To perform a brute-force scan of the sequence, use the `last!` function, keeping in mind that `last!` may never return a result.

#### An Example

```scheme
(define x '(99 64 32 48))
(last x)
```

This example will return _48_.

---
title: "fold-left"
description: "left folds a sequence"
names: ["reduce", "fold-left", "foldl"]
usage: "(fold-left func val? seq)"
tags: ["sequence", "comprehension"]
---

Iterates over a sequence, reducing its elements to a single resulting value. The function provided must take two arguments. The first and second sequence elements encountered are the initial values applied to that function. Thereafter, the result of the previous calculation is used as the first argument, while the next element is used as the second argument.

#### An Example

```scheme
(fold-left + 5 (range 1 11))
```

This will return the value _60_.

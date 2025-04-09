---
title: "partition"
description: "partitions a sequence"
names: ["partition"]
usage: "(partition count step? seq)"
tags: ["sequence", "comprehension"]
---

Partition a sequence into groups of _count_ elements, incrementing by the number of elements defined in _step_ (or _count_ if _step_ is not provided).

#### An Example

```scheme
(seq->list (partition 2 3 [1 2 3 4 5 6 7 8 9 10]))
```

This example will return _((1 2) (4 5) (7 8) (10))_.

---
title: "map"
description: "lazily maps sequences"
names: ["map"]
usage: "(map func seq+)"
tags: ["sequence"]
---

Creates a lazy sequence whose elements are the result of applying the provided function to the sequence elements. If more than one sequence is provided, their elements are retrieved in parallel to supply additional arguments to the mapped function. Mapping will terminate as soon as any sequence is exhausted.

#### An Example

```scheme
(map (lambda (x) (\* x 2)) [1 2 3 4])
```

This will return the lazy sequence _(2 4 6 8)_. The following example performs mapping in parallel.

```scheme
(map + [1 2 3 4] '(4 6 8) '(30 20 10 56))
```

This will return the lazy sequence _(35 28 21)_. Mapping only occurs three times in this example because the second sequence only has three elements.

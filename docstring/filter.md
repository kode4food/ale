---
title: "filter"
description: "lazily filters a sequence"
names: ["filter"]
usage: "(filter func seq)"
tags: ["sequence" "comprehension"]
---

Creates a lazy sequence whose content is the result of applying the provided function to the elements of the provided sequence. If the result of the application is truthy (not _#f_ (false) or the empty list) then the value will be included in the resulting sequence.

#### An Example

```scheme
(filter (lambda (x) (< x 3)) [1 2 3 4])
```

This will return the lazy sequence _(1 2)_

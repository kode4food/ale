---
title: "fold-right"
description: "reduces a sequence from right to left"
names: ["fold-right", "fold-right!", "foldr"]
usage: "(fold-right func val? seq) (fold-right! func val? seq)"
tags: ["sequence"]
---

These forms reduce a sequence from right to left. `fold-right` uses `reverse`, while `fold-right!` uses `reverse!` and can therefore work with non-reversible sequences by materializing them first. `foldr` is an alias for `fold-right`.

#### An Example

```scheme
(fold-right cons '() [1 2 3])
```

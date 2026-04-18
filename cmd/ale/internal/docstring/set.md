---
title: "set"
description: "creates a persistent set"
names: ["set"]
usage: "(set value*)"
tags: ["data", "sequence"]
---

Creates a persistent set containing unique values. Set literals can also be written with `#{...}` reader syntax.

#### An Example

```scheme
(set 1 2 2 3)
```

This example returns `#{1 2 3}`.

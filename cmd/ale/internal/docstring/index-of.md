---
title: "index-of"
description: "finds the first position of a value in a sequence"
names: ["index-of"]
usage: "(index-of coll value)"
tags: ["sequence"]
---

Returns the zero-based index of the first matching value in a sequence. If no value matches, returns _#f_.

#### An Example

```scheme
(index-of [4 8 15 16] 15)
```

This example returns `2`.

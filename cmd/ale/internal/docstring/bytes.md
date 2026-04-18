---
title: "bytes"
description: "creates a byte sequence"
names: ["bytes"]
usage: "(bytes value*)"
tags: ["data", "sequence"]
---

Creates a byte sequence from numeric values. Byte sequences also have reader syntax using `#b[...]`.

#### An Example

```scheme
(bytes 65 66 67)
```

This example returns the byte sequence representing `ABC`.

---
title: "append"
description: "appends a value to an appendable sequence"
names: ["append"]
usage: "(append seq value)"
tags: ["sequence"]
---

Returns a new sequence with `value` appended to the end of `seq`. The target must satisfy `appendable?`.

#### An Example

```scheme
(append [1 2 3] 4)
```

This example returns `[1 2 3 4]`.

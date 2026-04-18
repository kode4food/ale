---
title: "reverse"
description: "reverses a sequence"
names: ["reverse", "reverse!"]
usage: "(reverse coll) (reverse! coll)"
tags: ["sequence"]
---

`reverse` requires a reversible sequence. `reverse!` will reverse reversible sequences directly, and otherwise materialize the sequence before reversing it.

#### An Example

```scheme
(reverse [1 2 3 4])
```

This example returns `[4 3 2 1]`.

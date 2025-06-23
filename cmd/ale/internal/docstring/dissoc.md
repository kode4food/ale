---
title: "dissoc"
description: "removes an association from a mapped sequence by key"
names: ["dissoc", "dissoc*"]
usage: "(dissoc seq key) (dissoc* seq key+)"
tags: ["sequence"]
---

Returns a newly mapped sequence wherein the associations identified by the provided keys are removed. If the keys don't exist, the original sequence is returned.

#### An Example

```scheme
(define robert {:name "Bob" :age 45})
(dissoc robert :age)
```

This example returns a copy of _robert_ from which the _:age_ association has been removed. The original sequence is unaffected.

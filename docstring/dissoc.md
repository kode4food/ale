
---
title: "dissoc"
date: 2021-02-11T12:19:22+02:00
description: "removes an association by key"
names: ["dissoc"]
usage: "(dissoc seq key)"
tags: ["sequence"]
---

Returns a newly mapped sequence wherein the association identified by the key is removed. If the key doesn't exist, the original sequence is returned.

#### An Example

```scheme
(define robert {:name "Bob" :age 45})
(dissoc robert :age)
```

This example returns a copy of _robert_ from which the _:age_ association has been removed. The original sequence is unaffected.

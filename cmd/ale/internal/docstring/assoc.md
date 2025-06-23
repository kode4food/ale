---
title: "assoc"
description: "associates pairs with a mapped sequence"
names: ["assoc", "assoc*"]
usage: "(assoc seq pair) (assoc* seq pair+)"
tags: ["sequence"]
---

Returns a newly mapped sequence wherein the specified key/value pairs are associated. If a key already exists, the value replaces the one previously stored; otherwise the pair is added to the sequence.

#### An Example

```scheme
(define robert {:name "Bob" :age 45})
(assoc robert (:age . 46))
```

This example returns a copy of _robert_ wherein the value associated with _:age_ has been replaced by the number _46_. The original sequence is unaffected.

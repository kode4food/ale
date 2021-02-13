
---
title: "assoc"
date: 2021-02-11T12:19:22+02:00
description: "associates a value to a key"
names: ["assoc"]
usage: "(assoc seq key value) (assoc seq pair)"
tags: ["sequence"]
---

Returns a newly mapped sequence wherein the specified key and value are associated. If the key already exists, the value replaces the one previously stored, otherwise the pair are added to the sequence.

#### An Example

```scheme
(define robert {:name "Bob" :age 45})
(assoc robert :age 46)
```

This example returns a copy of _robert_ wherein the value associated with _:age_ has been replaced by the number _46_. The original object associated with _robert_ is unaffected.

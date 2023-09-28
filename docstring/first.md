---
title: "car"
description: "returns the first element of a pair or sequence"
names: ["car", "first"]
usage: "(car form)"
tags: ["sequence"]
---

This function will return the first element of the specified pair or sequence, or the empty list if the sequence is empty.

#### An Example

```scheme
(define x '(99 64 32 48))
(car x)  ;; will return 99

(define y (100 . 200))
(car y)  ;; will return 100
```

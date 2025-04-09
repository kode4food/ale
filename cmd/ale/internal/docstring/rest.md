---
title: "cdr"
description: "returns the rest of a pair or sequence"
names: ["cdr", "rest"]
usage: "(cdr seq)"
tags: ["sequence"]
---

This function will return the portion of a pair or sequence that excludes its first element. For sequences, this will be the remainder of the sequence. For cons pairs, this will be the cdr portion.

#### An Example

```scheme
(define x '(99 64 32 48))
(cdr x)  ;; will return (64, 32, 48)

(define y (100 . 200))
(cdr y)  ;; will return 200
```

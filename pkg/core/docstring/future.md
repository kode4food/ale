---
title: "future"
description: "evaluates a set of forms asynchronously"
names: ["future"]
usage: "(future form*)"
tags: ["concurrency"]
---

Returns a promise whose expressions are immediately evaluated in another thread of execution. 

#### An Example

```scheme
(define fut (future
  (seq->vector (generate
    (emit "red")
    (emit "orange")
    (emit "yellow")))))

(fut)
```

This example produces a promise called _fut_ that converts the results of an asynchronous block into a vector. The `(fut)` call will block until the future returns a value.

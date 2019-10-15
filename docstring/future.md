---
title: "future"
date: 2019-04-06T12:19:22+02:00
description: "evaluates a set of forms asynchronously"
names: ["future"]
usage: "(future form*)"
tags: ["concurrency"]
---
Returns a future in the form of a function. The provided forms will be evaluated in a separate thread of execution, and any calls to the function **will block** until the forms have been completely evaluated.

#### An Example

~~~scheme
(define fut (future
  (seq->vector (generate
    (emit "red")
    (emit "orange")
    (emit "yellow")))))

(fut)
~~~

This example produces a future called *fut* that converts the results of an asynchronous block into a vector. The `(fut)` call will block until the future returns a value.

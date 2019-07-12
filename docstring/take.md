---
title: "take"
description: "takes the first elements of a sequence"
date: 2019-04-06T12:19:22+02:00
names: ["take"]
usage: "(take count seq)"
tags: ["sequence", "comprehension"]
---
Will return a lazy sequence of either *count* or fewer elements from the beginning of the provided sequence. If the source sequence is shorter than the requested count, the resulting sequence will be truncated.

#### An Example

~~~scheme
(define x '(1 2 3 4))
(define y [5 6 7 8])
(take 6 (concat x y))
~~~

This example will return the lazy sequence _(1 2 3 4 5 6)_.

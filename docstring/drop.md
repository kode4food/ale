---
title: "drop"
date: 2019-04-06T12:19:22+02:00
description: "drops the first elements of a sequence"
names: ["drop"]
usage: "(drop count seq)"
tags: ["sequence", "comprehension"]
---
Will return a lazy sequence that excludes the first *count* elements of the provided sequence. If the source sequence is shorter than the requested count, an empty list will be returned.

#### An Example

~~~scheme
(def x '(1 2 3 4))
(def y [5 6 7 8])
(drop 3 (concat x y))
~~~

This example will return the lazy sequence _(4 5 6 7 8)_.

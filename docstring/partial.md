---
title: "partial"
date: 2019-04-06T12:19:22+02:00
description: "generates a function based on a partial apply"
names: ["partial"]
usage: "(partial func arg+)"
tags: ["function"]
---
Returns a new Function whose initial arguments are pre-bound to those provided. When that Function is invoked, any provided arguments will simply be appended to the pre-bound arguments before calling the original Function.

#### An Example

~~~scheme
(def plus10 (partial + 4 6))
(plus10 9)
~~~

This example will return _19_ as though `(+ 4 6 9)` were called.

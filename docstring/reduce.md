---
title: "reduce"
date: 2019-04-06T12:19:22+02:00
description: "reduces a sequence"
names: ["reduce"]
usage: "(reduce func val? seq)"
tags: ["sequence", "comprehension"]
---
Iterates over a set of sequence, reducing their elements to a single resulting value. The function provided must take two arguments. The first and second sequence elements encountered are the initial values applied to that function. Thereafter, the result of the previous calculation is used as the first argument, while the next element is used as the second argument.

#### An Example

~~~scheme
(reduce + 5 (range 1 11))
~~~
This will return the value _60_.

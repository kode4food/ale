---
title: "apply"
date: 2019-04-06T12:19:22+02:00
usage: "(apply func seq)"
names: ["apply"]
description: "applies a function to the provided arguments"
tags: ["function"]
---
Evaluates the provided sequence and applies the provided function to its values.

#### An Example

~~~scheme
(def x '(1 2 3))
(apply + x)
~~~

This example will return _6_.

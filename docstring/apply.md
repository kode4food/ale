---
title: "apply"
date: 2019-04-06T12:19:22+02:00
usage: "(apply func seq)"
names: ["apply"]
description: "applies arguments to a function"
tags: ["function"]
---
Evaluates the provided sequence and applies its values to the provided function (or applicable) as its arguments.

#### An Example

```clojure
(def x '(1 2 3))
(apply + x)
```

This example will return _6_.

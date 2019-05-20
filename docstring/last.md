---
title: "last"
date: 2019-04-06T12:19:22+02:00
description: "returns the last element of the sequence"
names: ["last"]
usage: "(last seq) (last! seq)"
tags: ["sequence"]
---
This function will return the last element of the specified sequence, or _nil_ if the sequence is empty. If the sequence is lazily computed, asynchronous, or otherwise incapable of returning a count, this function will raise an error.

In order to perform a brute-force scan of the sequence, use the `last!` function, keeping in mind that `last!` may never return a result.

#### An Example

```clojure
(def x '(99 64 32 48))
(last x)
```

This example will return _48_.

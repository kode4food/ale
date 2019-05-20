---
title: "last"
date: 2019-04-06T12:19:22+02:00
description: "returns the last element of the sequence"
names: ["last"]
usage: "(last seq)"
tags: ["sequence"]
---
This function will return the last element of the specified sequence, or _nil_ if the sequence is empty.

#### An Example

```clojure
(def x '(99 64 32 48))
(last x)
```

This example will return _48_.

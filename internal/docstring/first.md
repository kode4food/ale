---
title: "first"
date: 2019-04-06T12:19:22+02:00
description: "returns the first element of the sequence"
names: ["first"]
usage: "(first seq)"
tags: ["sequence"]
---
This function will return the first element of the specified sequence, or _nil_ if the sequence is empty. It would be beneficial to check for a valid sequence using `(seq? seq)` before calling `first` or `rest`.

#### An Example

```clojure
(def x '(99 64 32 48))
(first x)
```

 This example will return _99_.

---
title: "to-list"
date: 2019-04-06T12:19:22+02:00
description: "converts sequences to a list"
names: ["to-list"]
usage: "(to-list seq+)"
tags: ["sequence", "conversion"]
---
Will concatenate a set of sequences into a list. Unlike the standard `concat` function, which is lazily computed, the result of this function will be materialized immediately.

#### An Example

```clojure
(def x [1 2 3 4])
(def y
  (map (fn [x] (+ x 4))
  '(1 2 3 4)))
(to-list x y)
```

---
title: "to-assoc"
date: 2019-04-06T12:19:22+02:00
description: "converts sequences to an associative structure"
names: ["to-assoc"]
usage: "(to-assoc seq+)"
tags: ["sequence", "conversion"]
---
Will concatenate a set of sequences into an associative. Unlike the standard `concat` function, which is lazily computed, the result of this function will be materialized immediately.

#### An Example

```clojure
(def x [:name "ale" :age 0.3])
(def y '(:weight "light"))
(to-assoc x y)
```

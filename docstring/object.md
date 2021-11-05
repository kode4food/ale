---
title: "object"
description: "creates a new object instance"
names: ["object"]
usage: "(object <key value>*)"
tags: ["data"]
---

Will create a new object (hash-map) based on the provided key-value pairs, or return an empty object if no forms are provided. This function is no different from the object literal syntax except that it can be treated in a first-class fashion.

An object can be iterated over as a sequence. The resulting sequence is guaranteed to have no duplicated keys, but is not guaranteed to return in any particular order.

#### An Example

```scheme
(object :name "ale" :age 0.3 :lang "Go")

;; which is the same as
{:name "ale" :age 0.3 :lang "Go"}
```

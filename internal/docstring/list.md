---
title: "list"
date: 2019-04-06T12:19:22+02:00
description: "creates a new list"
names: ["list"]
usage: "(list form*)"
tags: ["sequence"]
---
Will create a new list whose elements are the evaluated forms provided, or return the empty list if no forms are provided.

#### An Example

```clojure
(def x "hello")
(def y "there")
(list x y)
```

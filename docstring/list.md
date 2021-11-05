---
title: "list"
description: "creates a new list"
names: ["list"]
usage: "(list form*)"
tags: ["sequence"]
---

Will create a new list whose elements are the evaluated forms provided, or return the empty list if no forms are provided.

#### An Example

```scheme
(define x "hello")
(define y "there")
(list x y)
```

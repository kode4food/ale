---
title: "vector"
description: "creates a new vector"
names: ["vector"]
usage: "(vector form*)"
tags: ["sequence"]
---

Will create a new vector whose elements are the evaluated forms provided. This function is no different from the vector literal syntax except that it can be treated in a first-class fashion.

#### An Example

```scheme
(define x "hello")
(define y "there")
(vector x y)
```

---
title: "juxt"
description: "juxtaposes a set of functions"
names: ["juxt"]
usage: "(juxt func*)"
tags: ["function"]
---

Returns a new function that represents the juxtaposition of the provided functions. This function returns a vector containing the result of applying each provided function to the juxtaposed function's arguments.

#### An Example

```scheme
(define juxt-math (juxt \* + - /))
(juxt-math 32 10)
```

This example will return _[320 42 22 16/5]_ as though `[(\* 32 10) (+ 32 10) (- 32 10) (/ 32 10)]` were called.

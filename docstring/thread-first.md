---
title: "thread first (->)"
description: "threads value through calls as their first argument"

names: ["->"]
usage: "(-> expr forms*)"
tags: ["function"]
---

Evaluates _expr_ and threads it through the supplied forms as their first argument. Any form that is not already a function call will be converted into one before threading.

#### An Example

```scheme
(-> 0 (+ 10) (\* 2) (/ 5))
```

Will expand to `(/ (\* (+ 0 10) 2) 5)` and return _4_.

---
title: "thread first (->)"
description: "threads value through calls as their first argument"
date: 2019-04-06T12:19:22+02:00
names: ["->"]
usage: "(-> expr forms*)"
tags: ["function"]
---
Evaluates *expr* and threads it through the supplied forms as their first argument. Any form that is not already a function call will be converted into one before threading.

#### An Example

```clojure
(-> 0 (+ 10) (* 2) (/ 5))
```

Will expand to `(/ (* (+ 0 10) 2) 5)` and return _4_.

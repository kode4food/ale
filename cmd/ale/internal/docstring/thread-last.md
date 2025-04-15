---
title: "thread last (->>)"
description: "threads value through calls as their last argument"
names: ["->>"]
usage: "(->> expr forms*)"
tags: ["function"]
---

Evaluates _expr_ and threads it through the supplied forms as their last argument. Any form that is not already a function call will be converted into one before threading.

#### An Example

```scheme
(->> 0 (+ 10) (\* 2) (/ 5.0))
```

Will expand to `(/ 5.0 (\* 2 (+ 10 0)))` and return _0.25_.

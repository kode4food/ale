---
title: "delay"
date: 2019-04-06T12:19:22+02:00
description: "produces a delayed evaluation"
names: ["delay"]
usage: "(delay expr*)"
tags: ["concurrency"]
---

Returns a promise that, when forced, evaluates the exprs, returning the final evaluated result. The result is then cached, so further uses of force return the cached value immediately.

#### An Example

```scheme
(define p (delay
              (println "hello once")
              "hello"))
(force p) ;; prints "hello once"
(force p)
```

The first invocation of `p` will print _"hello once"_ to the console, and also return the string _"hello"_. Subsequent invocations of `p` will only return _"hello"_.

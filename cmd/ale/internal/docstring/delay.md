---
title: "delay"
description: "produces a delayed evaluation"
names: ["delay", "force", "force!", "delay-force"]
usage: "(delay expr*) (force value) (force! value) (delay-force expr*)"
tags: ["concurrency"]
---

`delay` returns a promise that evaluates its body the first time it is forced. That result is cached, so later calls return immediately. `force` resolves one promise layer, or returns non-promises unchanged. `force!` keeps forcing until the result is no longer a promise. `delay-force` delays a computation whose result should be forced once before being cached.

#### An Example

```scheme
(define p (delay
            (println "hello once")
            "hello"))
(force p) ;; prints "hello once"
(force p)
```

The first invocation of `p` will print _"hello once"_ to the console, and also return the string _"hello"_. Subsequent invocations of `p` will only return _"hello"_.

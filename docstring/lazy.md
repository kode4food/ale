---
title: "lazy"
date: 2019-04-06T12:19:22+02:00
description: "produces a lazy evaluation"
names: ["lazy"]
usage: "(lazy expr*)"
tags: ["concurrency"]
---

Like `delay`, but if the initial forced result is a promise, it will continue to be forced until a non-promise result is capable of being returned.

#### An Example

```scheme
(define p (lazy
              (println "hello once")
              (delay
                (println "hello twice")
                "hello")))
(force p) ;; prints "hello once / hello twice"
(force p)
```

The first invocation of `p` will print _"hello once"_ followed by _"hello twice"_ to the console, also returning the string _"hello"_. Subsequent invocations of `p` will only return _"hello"_.

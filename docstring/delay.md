---
title: "delay"
date: 2019-04-06T12:19:22+02:00
description: "produces a lazy evaluation"
names: ["delay"]
usage: "(delay expr)"
tags: ["concurrency"]
---
Returns a promise in the form of a function. The provided expression will only be evaluated when the promise function is invoked. The expression will only be evaluated once, and then memoized.

#### An Example

~~~scheme
(define p (delay
            (begin
              (println "hello once")
              "hello")))
(p) ;; prints "hello once"
(p)
~~~

The first invocation of `p` will print _"hello once"_ to the console, and also return the string _"hello"_. Subsequent invocations of `p` will only return _"hello"_.
---
title: "go"
date: 2019-04-06T12:19:22+02:00
description: "asynchronously evaluates a block"
names: ["go"]
usage: "(go form*)"
tags: ["concurrency"]
---

The provided forms will be evaluated in a separate thread of execution. Any resulting value of the block will be discarded.

#### An Example

```scheme
(define ch (chan))
(go (: ch :emit "hello")
    (: ch :close))
(str (first (ch :seq)) " world!")
```

---
title: "promise"
date: 2019-04-06T12:19:22+02:00
description: "produces an unresolved value"
names: ["promise"]
usage: "(promise value?)"
tags: ["concurrency"]
---
Returns a promise in the form of a function. If applied without an argument, this function **will block**, waiting for a value to be delivered via a call to the promise function that includes an argument. A promise can only be delivered once.

#### An Example

```clojure
(def p (promise))
(p "hello")
(p)
```

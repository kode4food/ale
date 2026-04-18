---
title: "letfn"
description: "binds mutually recursive local functions"
names: ["letfn"]
usage: "(letfn [binding*] form*)"
tags: ["binding", "macro"]
---

Creates a local scope containing named recursive function bindings. It is useful when several local procedures need to reference one another.

#### An Example

```scheme
(letfn [(lambda-rec even? (n) (if (zero? n) true (odd? (dec n))))
        (lambda-rec odd?  (n) (if (zero? n) false (even? (dec n))))]
  (even? 10))
```

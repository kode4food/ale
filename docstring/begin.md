---
title: "begin"
description: "evaluates a sequence of forms"
names: ["begin"]
usage: "(begin form*)"
---

Will evaluate each form in turn, returning the final evaluation as its result.

#### An Example

```scheme
(begin
  (println "hello")
  "returned")
```

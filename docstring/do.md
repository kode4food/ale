---
title: "do"
date: 2019-04-06T12:19:22+02:00
description: "evaluates multiple forms"
names: ["do"]
usage: "(do form*)"
---
Will evaluate each form in turn, returning the final evaluation as its result.

#### An Example

```clojure
(do
  (println "hello")
  "returned")
```

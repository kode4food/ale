---
title: "printing"
description: "writes values to standard output"
names: ["print", "println", "pr", "prn"]
usage: "(print form*) (println form*) (pr form*) (prn form*)"
tags: ["io"]
---

These functions print space-separated values to `*out*`. `print` and `println` use the regular string conversion. `pr` and `prn` use reader-oriented output via `str!`. `println` and `prn` append a trailing newline.

#### An Example

```scheme
(println "hello" 42)
(prn {:name "ale"})
```

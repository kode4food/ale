---
title: "branching helpers"
description: "conditional convenience macros"
names: ["case", "if-let", "when-let"]
usage: "(case expr [test branch]*) (if-let [name expr] then else?) (when-let [name expr] form*)"
tags: ["branching", "macro"]
---

These macros build on `if` to cover common control-flow patterns. `case` compares a value using `eq`, and a clause test may be either a single value or a list of values. `if-let` binds a value and tests whether it is truthy. `when-let` is the body-only form of `if-let`. If no `case` clause matches, it raises an error.

#### An Example

```scheme
(if-let [x (get {:name "ale"} :title false)]
  (str "hello " x)
  "missing")
```

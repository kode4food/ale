---
title: "type helpers"
description: "type queries and assertions"
names: ["assert-type", "is-a"]
usage: "(assert-type type value) (is-a type value)"
tags: ["type"]
---

These helpers work with Ale types. `assert-type` returns `value` when it matches `type`, otherwise it raises. `is-a` looks up the predicate for a builtin type keyword and applies it to `value`.

#### An Example

```scheme
(assert-type :number 42)
(is-a :string "hello")
```

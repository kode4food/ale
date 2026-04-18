---
title: "raise"
description: "raises an error value"
names: ["raise"]
usage: "(raise form*)"
tags: ["exception"]
---

Converts its arguments to a string and raises the resulting value as an error.

#### An Example

```scheme
(raise "unexpected value: " 42)
```

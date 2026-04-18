---
title: "syntax-quote"
description: "quotes code while preserving symbol hygiene"
names: ["syntax-quote"]
usage: "(syntax-quote form)"
tags: ["macro"]
---

Quotes a form while resolving symbols hygienically. Within a syntax-quoted form, `unquote` and `unquote-splicing` can inject evaluated values.

#### An Example

```scheme
(let [x 99]
  `(list ,x))
```

---
title: "gensym"
description: "creates a unique symbol, useful in macros"
names: ["gensym"]
usage: "(gensym sym?)"
tags: ["symbol", "macro"]
---

If an unqualified symbol is provided, that symbol will be used to clarify the uniquely generated symbol. This function provides the underlying behavior for hash-tailed symbols in syntax-highlighting macros.

#### An Example

```scheme
(let [s (gensym 'var)]
  (list 'ale/let [s "hello"] s))

;; is equivalent to
``(let [var# "hello"] var#)
```

This example will return _"hello"_.

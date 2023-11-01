---
title: "not"
description: "logically inverts the truthiness of the provided form"
names: ["not"]
usage: "(not form)"
tags: ["logic"]
---

Return _#f_ (false) if the provided _form_ is truthy, otherwise will return _#t_ (true).

#### An Example

```scheme
(not "hello")
```

This will return the boolean _#f_ (false) because the value _"hello"_ is truthy.

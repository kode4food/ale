---
title: "seq"
description: "attempts to convert form to a sequence"
names: ["seq"]
usage: "(seq form)"
tags: ["sequence" "conversion"]
---

Will attempt to convert the provided form to a sequence if it isn't already. If the form cannot be converted, or if the resulting sequence is empty, the empty list will be returned.

#### An Example

```scheme
(when-let [s (seq "hello")]
  (seq->vector (map (lambda (x) (str x "-")) s)))
```

This example will return _["h-" "e-" "l-" "l-" "o-"]_.

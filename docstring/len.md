---
title: "len"
date: 2019-04-06T12:19:22+02:00
description: "returns the size of a sequence"
names: ["len", "len!"]
usage: "(len seq) (len! seq)"
tags: ["sequence"]
---
The `len` function will return the number of elements in a sequence. If the sequence is lazily computed, asynchronous, or otherwise incapable of returning a count, this function will raise an error.

In order to perform a brute-force count of a sequence, use the `len!` function, keeping in mind that `len!` may never return a result.

#### An Example

```clojure
(len '(1 2 3 4))
```

This example will return _4_.

---
title: "len"
date: 2019-04-06T12:19:22+02:00
description: "returns the length of a sequence"
names: ["len"]
usage: "(len seq)"
tags: ["sequence"]
---
Will return the number of elements in a countable sequence. If the sequence is lazily computed, asynchronous, or otherwise incapable of being counted, this function will raise an error.

#### An Example

```clojure
(len '(1 2 3 4))
```

This example will return _4_.

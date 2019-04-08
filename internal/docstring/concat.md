---
title: "concat"
date: 2019-04-06T12:19:22+02:00
description: "lazily concatenates sequences"
names: ["concat"]
usage: "(concat seq+)"
tags: ["sequence", "comprehension"]
---
Creates a lazy sequence whose content is the result of concatenating the elements of each provided sequence.

#### An Example

```clojure
(to-list (concat [1 2 3] '(4 5 6)))
```

This will return the list _(1 2 3 4 5 6)_

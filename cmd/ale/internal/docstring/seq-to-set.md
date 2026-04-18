---
title: "seq->set"
description: "concatenates sequences into a set"
names: ["seq->set"]
usage: "(seq->set seq+)"
tags: ["sequence", "data"]
---

Concatenates one or more sequences into a set. Duplicate values collapse to a single member.

#### An Example

```scheme
(seq->set [1 2] '(2 3))
```

This example returns `#{1 2 3}`.

---
title: "bytes?"
description: "tests whether the provided forms are byte sequences"
names: ["bytes?", "!bytes?"]
usage: "(bytes? form+) (!bytes? form+)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to byte sequences, this function returns _#t_ (true). Otherwise it returns _#f_ (false).

#### An Example

```scheme
(bytes? (bytes 65 66) #b[1 2 3])
```

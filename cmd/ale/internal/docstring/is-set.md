---
title: "set?"
description: "tests whether the provided forms are sets"
names: ["set?", "!set?"]
usage: "(set? form+) (!set? form+)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to sets, this function returns _#t_ (true). Otherwise it returns _#f_ (false).

#### An Example

```scheme
(set? #{1 2 3} (set 4 5 6))
```

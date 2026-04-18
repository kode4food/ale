---
title: "contains?"
description: "tests whether a mapped value contains a key"
names: ["contains?"]
usage: "(contains? coll key)"
tags: ["predicate", "sequence"]
---

Returns _#t_ if `coll` can resolve `key`, otherwise _#f_. This works with mapped lookup types such as objects and sets.

#### An Example

```scheme
(contains? #{:name :age} :name)
```

This example returns _#t_.

---
title: "mapped?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are mapped"
names: ["mapped?" "!mapped?" "is-mapped"]
usage: "(mapped? form+) (!mapped? form+) (is-mapped form)"
tags: ["sequence" "predicate"]
---

If all forms evaluate to a mapped type, then this function will return _#t_ (true). The first non-mapped will result in the function returning _#f_ (false).

#### An Example

```scheme
(mapped? {:name "bill"} {:name "peggy"} [1 2 3])
```

This example will return _#f_ (false) because the third form is a vector.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be mapped.

```scheme
(!mapped? "hello" [1 2 3])
```

This example will return _#t_ (true).

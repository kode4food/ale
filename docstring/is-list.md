---
title: "list?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are lists"
names: ["list?", "!list?", "is-list"]
usage: "(list? form+) (!list? form+) (is-list form)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to a list, then this function will return _#t_ (true). The first non-list will result in the function returning _#f_ (false).

#### An Example

```scheme
(list? '(1 2 3 4) [5 6 7 8])
```

This example will return _#f_ (false) because the second form is a vector.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be lists.

```scheme
(!list? "hello" [99])
```

This example will return _#t_ (true).

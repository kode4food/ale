---
title: "null?"
description: "tests whether the provided forms are nil"
names: ["null?", "!null?", "is-null"]
usage: "(null? form+) (!null? form+) (is-null form)"
tags: ["predicate"]
---

If all forms evaluate to null (empty list), then this function will return _#t_ (true). The first non-null will result in the function returning _#f_ (false).

#### An Example

```scheme
(null? '(1 2 3 4) '())
```

This example will return _#f_ (false) because the first form is a list.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be nil.

```scheme
(!null? "hello" [99])
```

This example will return _#t_ (true).

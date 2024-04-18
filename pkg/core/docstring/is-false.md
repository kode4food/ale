---
title: "false?"
description: "tests whether the provided forms are boolean false"
names: ["false?", "!false?", "is-false"]
usage: "(false? form+) (!false? form+) (is-false form)"
tags: ["predicate"]
---

If all forms evaluate to false (_#f_), then this function will return _#t_ (true). The first non-false will result in the function returning _#f_ (false).

#### An Example

```scheme
(false? (< 3 2) (> 5 10))
```

This example will return _#t_ (true) because all the equalities result in _#f_ (false).

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be true. This is equivalent to the _true?_ predicate.

```scheme
(!false? #t (< 5 10))
```

This example will return _#t_ (true).

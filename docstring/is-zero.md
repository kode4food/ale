---
title: "zero?"
description: "tests whether the provided forms are numeric zero"
names: ["zero?", "!zero?", "is-zero"]
usage: "(zero? form+) (!zero? form+) (is-zero form)"
tags: ["predicate"]
---

If all forms evaluate to zero (0, 0.0, or 0/x), then this function will return _#t_ (true). The first non-zero will result in the function returning _#f_ (false).

#### An Example

```scheme
(zero? (- 3 2 1) 0 0/1 (/ 3 3))
```

This example will return _#f_ (false) because the last form evaluates to _1_.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be zero.

```scheme
(!zero? "hello" [99])
```

This example will return _#t_ (true).

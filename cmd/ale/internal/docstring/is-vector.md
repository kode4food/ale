---
title: "vector?"
description: "tests whether the provided forms are vectors"
names: ["vector?", "!vector?", "is-vector"]
usage: "(vector? form+) (!vector? form+) (is-vector form)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to a vector, then this function will return _#t_ (true). The first non-vector will result in the function returning _#f_ (false).

#### An Example

```scheme
(vector? '(1 2 3 4) [5 6 7 8])
```

This example will return _#f_ (false) because the first form is a list.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be vectors.

```scheme
(!vector? "hello" [99])
```

This example will return _#f_ (false) because the second form is a vector.

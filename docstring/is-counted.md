---
title: "counted?"
description: "tests whether the provided forms are counted sequences"
names: ["counted?", "!counted?", "is-counted"]
usage: "(counted? form+) (!counted? form+) (is-counted form)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to a valid sequence than can report its length without counting, then this function will return _#t_ (true). The first non-counted sequence will result in the function returning _#f_ (false).

#### An Example

```scheme
(counted? '(1 2 3 4) [5 6 7 8])
```

This example will return _#t_ (true).

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be valid counted sequences.

```scheme
(!counted? "hello" 99)
```

This example will return _#f_ (false) because strings are counted sequences.

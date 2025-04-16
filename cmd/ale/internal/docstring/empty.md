---
title: "empty?"
description: "tests whether the provided forms are empty sequences"
names: ["empty?", "!empty?"]
usage: "(empty? form+) (!empty? form+)"
tags: ["predicate"]
---

If all forms evaluate to empty sequences, then this function will return _#t_ (true). The first evaluation that is not an empty sequence will result in the function returning _#f_ (false).

#### An Example

```scheme
(empty? '(1 2 3 4) [] {})
```

This example will return _#f_ (false) because the first form is a list with four elements.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be nil.

```scheme
(!empty? '(1) [2] "3" {4 5})
```

This example will return _#t_ (true) because all the arguments are non-empty sequences of some kind.

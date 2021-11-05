---
title: "atom?"
description: "tests whether the provided forms are atomic"
names: ["atom?" "!atom?" "is-atom"]
usage: "(atom? form+) (!atom? form+) (is-atom form)"
tags: ["predicate"]
---

A form is considered to be atomic if it cannot be further evaluated and would otherwise evaluate to itself.

#### An Example

```scheme
(atom? '() :hello "there")
```

This example will return _#t_ (true) because each value is atomic.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be atomic.

```scheme
(!atom? '(+ 1 2 3) [4 5 6])
```

This example will return _#t_ (true) because compound types such as lists and vectors are not considered to be atomic.

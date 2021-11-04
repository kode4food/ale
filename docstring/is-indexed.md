---
title: "indexed?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are indexed sequences"
names: ["indexed?" "!indexed?" "is-indexed"]
usage: "(indexed? form+) (!indexed? form+) (is-indexed form)"
tags: ["sequence" "predicate"]
---

If all forms evaluate to a valid sequence than can be accessed by index, then this function will return _#t_ (true). The first non-indexed sequence will result in the function returning _#f_ (false).

#### An Example

```scheme
(indexed? '(1 2 3 4) [5 6 7 8])
```

This example will return _#t_ (true).

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be valid indexed sequences.

```scheme
(!indexed? "hello" 99)
```

This example will return _#f_ (false) because strings are indexed sequences.

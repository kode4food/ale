---
title: "sized?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are sized sequences"
names: ["sized?", "!sized?", "is-sized"]
usage: "(sized? form+) (!sized? form+) (is-sized form)"
tags: ["sequence", "predicate"]
---
If all forms evaluate to a valid sequence than can report its length without counting, then this function will return _true_. The first non-sized sequence will result in the function returning _false_.

#### An Example

```clojure
(sized? '(1 2 3 4) [5 6 7 8])
```

This example will return _true_.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be valid sized sequences.

```clojure
(!sized? "hello" 99)
```

This example will return _false_ because strings are sized sequences.

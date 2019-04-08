---
title: "seq?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are sequences"
names: ["seq?", "!seq?", "is-seq"]
usage: "(seq? form+) (!seq? form+) (is-seq form)"
tags: ["sequence", "predicate"]
---
If all forms evaluate to a valid sequence, then this function will return _true_. The first non-sequence will result in the function returning _false_.

#### An Example

```clojure
(seq? '(1 2 3 4) [5 6 7 8])
```

This example will return _true_.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be valid sequences.

```clojure
(!seq? "hello" 99)
```

This example will return _true_.

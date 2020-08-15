---
title: "promise?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are promises"
names: ["promise?", "!promise?", "is-promise"]
usage: "(promise? form+) (!promise? form+) (is-promise form)"
tags: ["concurrency", "predicate"]
---

If all forms evaluate to a promise, then this function will return _#t_ (true). The first non-promise will result in the function returning _#f_ (false).

#### An Example

```scheme
(define p1 (promise "one"))
(define p2 (promise "two"))
(promise? p1 p2 [1 2 3])
```

This example will return _#f_ (false) because the third form is a vector.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be promises.

```scheme
(!promise? "hello" [1 2 3])
```

This example will return _#t_ (true).

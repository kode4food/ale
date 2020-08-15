---
title: "object?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are objects"
names: ["object?", "!object?", "is-object"]
usage: "(object? form+) (!object? form+) (is-object form)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to an object (hash-map), then this function will return _#t_ (true). The first non-object will result in the function returning _#f_ (false).

#### An Example

```scheme
(object? {:name "bill"} {:name "peggy"} [1 2 3])
```

This example will return _#f_ (false) because the third form is a vector.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be objects.

```scheme
(!object? "hello" [1 2 3])
```

This example will return _#t_ (true).

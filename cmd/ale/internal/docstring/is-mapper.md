---
title: "mapper?"
description: "tests whether the provided forms support association"
names: ["mapper?", "!mapper?"]
usage: "(mapper? form+) (!mapper? form+)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to mapper values, then this function will return _#t_ (true). A mapper supports keyed lookup and association operations such as `assoc` and `dissoc`.

#### An Example

```scheme
(mapper? {:name "bill"} #{:name :age} [1 2 3])
```

This example will return _#f_ (false) because sets and vectors are not mappers.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be mappers.

```scheme
(!mapper? #{:name :age} [1 2 3])
```

This example will return _#t_ (true).

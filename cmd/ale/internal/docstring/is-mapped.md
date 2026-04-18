---
title: "mapped?"
description: "tests whether the provided forms support keyed lookup"
names: ["mapped?", "!mapped?"]
usage: "(mapped? form+) (!mapped? form+)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to mapped values, then this function will return _#t_ (true). A mapped value supports keyed lookup through functions like `get` and `contains?`.

#### An Example

```scheme
(mapped? {:name "bill"} #{:name :age} [1 2 3])
```

This example will return _#f_ (false) because the third form is a vector.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be mapped.

```scheme
(!mapped? "hello" [1 2 3])
```

This example will return _#t_ (true).

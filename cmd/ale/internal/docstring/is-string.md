---
title: "str?"
description: "tests whether the provided forms are strings"
names: ["string?", "!string?"]
usage: "(string? form+) (!string? form+)"
tags: ["sequence", "predicate"]
---

If all forms evaluate to string, then this function will return _#t_ (true). The first non-string will result in the function returning _#f_ (false).

#### An Example

```scheme
(string? '(1 2 3 4) "hello")
```

This example will return _#f_ (false) because the first form is a list.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be strings.

```scheme
(!string? '(1 2 3) [99])
```

This example will return _#t_ (true).

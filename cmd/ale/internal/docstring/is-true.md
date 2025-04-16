---
title: "true?"
description: "tests whether the provided forms are boolean true"
names: ["true?", "!true?"]
usage: "(true? form+) (!true? form+)"
tags: ["predicate"]
---

If all forms evaluate to true (_#t_), then this function will return _#t_ (true). The first non-true will result in the function returning _#f_ (false). 

#### An Example

```scheme
(true? (> 3 2) (< 5 10) (and #t #f))
```

This example will return _#f_ (false) because the last form evaluates to _#f_.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be true. This is equivalent to the _false?_ predicate.

```scheme
(!true? #f (> 5 10))
```

This example will return _#t_ (true).

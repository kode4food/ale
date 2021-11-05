---
title: "identical (eq)"
description: "tests if a set of values are identical to the first"
names: ["eq" "!eq"]
usage: "(eq form form+) (!eq form form+)"
tags: ["comparison"]
---

Will return _#f_ (false) as soon as it encounters a form that is not identical to the first. Otherwise will return _#t_ (true).

#### An Example

```scheme
(define h "hello")
(eq "hello" h)
```

Like most predicates, this function can also be negated by prepending the `!` character. In this case, _#t_ (true) will be returned if not all forms are equal.

---
title: "list?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are lists"
names: ["list?", "!list?", "is-list"]
usage: "(list? form+) (!list? form+) (is-list form)"
tags: ["sequence", "predicate"]
---
If all forms evaluate to a list, then this function will return _true_. The first non-list will result in the function returning _false_.

#### An Example

```clojure
(list? '(1 2 3 4) [5 6 7 8])
```

This example will return _false_ because the second form is a vector.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be lists.

```clojure
(!list? "hello" [99])
```

This example will return _true_.

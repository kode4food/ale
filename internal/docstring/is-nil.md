---
title: "nil?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are nil"
names: ["nil?", "!nil?", "is-nil"]
usage: "(nil? form+) (!nil? form+) (is-nil form)"
tags: ["predicate"]
---
If all forms evaluate to nil, then this function will return _true_. The first non-nil will result in the function returning _false_.

#### An Example

```clojure
(nil? '(1 2 3 4) nil)
```

This example will return _false_ because the first form is a list.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be nil.

```clojure
(!nil? "hello" [99])
```

This example will return _true_.

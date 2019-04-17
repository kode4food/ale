---
title: "atom?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are atomic"
names: ["atom?", "!atom?", "is-atom"]
usage: "(atom? form+) (!atom? form+) (is-atom form)"
tags: ["predicate"]
---
A form is considered to be atomic if it cannot be further evaluated and would otherwise evaluate to itself.

#### An Example

```clojure
(atom? nil :hello "there")
```

This example will return _true_ because each value is atomic.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be atomic.

```clojure
(!atom? '(+ 1 2 3) [4 5 6])
```

This example will return _true_ because compound types such as lists and vectors are not considered to be atomic.

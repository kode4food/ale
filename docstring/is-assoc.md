---
title: "assoc?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are associatives"
names: ["assoc?", "!assoc?", "is-assoc"]
usage: "(assoc? form+) (!assoc? form+) (is-assoc form)"
tags: ["sequence", "predicate"]
---
If all forms evaluate to an assoc, then this function will return _#t_ (true). The first non-assoc will result in the function returning _#f_ (false).

#### An Example

~~~scheme
(assoc? {:name "bill"} {:name "peggy"} [1 2 3])
~~~

This example will return _#f_ (false) because the third form is a vector.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be associatives.

~~~scheme
(!assoc? "hello" [1 2 3])
~~~

This example will return _#t_ (true).

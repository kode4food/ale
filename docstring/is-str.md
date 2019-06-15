---
title: "str?"
date: 2019-04-06T12:19:22+02:00
description: "tests whether the provided forms are strings"
names: ["str?", "!str?", "is-str"]
usage: "(str? form+) (!str? form+) (is-str form)"
tags: ["sequence", "predicate"]
---
If all forms evaluate to strings, then this function will return _#t_ (true). The first non-string will result in the function returning _#f_ (false).

#### An Example

~~~scheme
(str? '(1 2 3 4) "hello")
~~~

This example will return _#f_ (false) because the first form is a list.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all of the provided forms must not be strings.

~~~scheme
(!str? '(1 2 3) [99])
~~~

This example will return _#t_ (true).

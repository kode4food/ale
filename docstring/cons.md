---
title: "cons"
date: 2019-04-06T12:19:22+02:00
description: "combines an element with a sequence"
names: ["cons"]
usage: "(cons form seq)"
tags: ["sequence"]
---
With an ordered sequence, such as a list or vector, the result is a new list or vector with the form prepended to the original. With an unordered sequence, such as an associative array, there is no guarantee regarding position.

The name *cons* is a vestige of when Lisp implementations **cons**tructed new lists or cells by pairing a _car_ (**c**ontents of the **a**ddress part of **r**egister) with a _cdr_ (**c**ontents of the **d**ecrement part of **r**egister).

#### An Example

```clojure
(def x '(3 4 5 6))
(def y (cons 2 x))
(cons 1 y)
```

This example will return _(1 2 3 4 5 6)_.

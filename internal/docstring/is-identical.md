---
title: "identical (eq)"
date: 2019-04-06T12:19:22+02:00
description: "tests if a set of values are identical to the first"
names: ["eq", "!eq", "is-eq"]
usage: "(eq form form+) (!eq form form+) (is-eq form form)"
tags: ["comparison"]
---
Will return _false_ as soon as it encounters a form that is not identical to the first. Otherwise will return _true_.

#### An Example

```clojure
(def h "hello")
(eq "hello" h)
```

Like most predicates, this function can also be negated by prepending the `!` character. In this case, _true_ will be returned if not all forms are equal.

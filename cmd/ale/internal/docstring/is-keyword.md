---
title: "keyword?"
description: "tests whether the provided forms are keywords"
names: ["keyword?", "!keyword?"]
usage: "(keyword? form+) (!keyword? form+)"
tags: ["predicate"]
---

A form is considered to be a keyword if it is a symbol with a colon (":") prefix. Keywords are often used as unique identifiers and are typically used to access values in objects or define objects keys.

#### An Example

```scheme
(keyword? :keyword1 :keyword2 :keyword3)
```

This example will return _#t_ (true) because each provided form is a keyword.

Like most predicates, this function can also be negated by prepending the `!` character. This means that all the provided forms must not be keywords.

```scheme
(!keyword? "string1" "string2" :keyword1)
```

This example will return _#f_ (false) because _:keyword1_ is provided.

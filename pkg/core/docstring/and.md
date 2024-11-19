---
title: "and"
description: "performs a short-circuiting and"
names: ["and"]
usage: "(and form*)"
tags: ["conditional", "logic"]
---

Evaluates forms from left to right. As soon as one evaluates to false, will return that value. Otherwise, it will proceed to evaluate the next form.

#### An Example

```scheme
(and (+ 1 2 3)
     false
     "not returned")
```

Will return _#f_ (false), never evaluating _"not returned"_, whereas:

```scheme
(and (+ 1 2 3)
     true
     "returned")
```

Will return the string _"returned"_.

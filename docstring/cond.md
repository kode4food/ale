---
title: "cond"
description: "performs conditional branching"
names: ["cond"]
usage: "(cond [pred then]*)"
tags: ["conditional"]
---

For each _pred-then_ clause, the predicate will be evaluated, and if it is truthy (not _#f_ (false) or the empty list) the _then_ form is evaluated and returned, otherwise the next clause is processed.

#### An Example

```scheme
(define x 99)

(cond
  [(< x 50)  "was less than 50"    ]
  [(> x 100) "was greater than 100"]
  [:else     "was in between"      ])
```

In this case, _"was in between"_ will be returned. The reason that this works is that the `:else` keyword, like all keywords, evaluates to truthy.

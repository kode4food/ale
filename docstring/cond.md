---
title: "cond"
date: 2019-04-06T12:19:22+02:00
description: "performs conditional branching"
names: ["cond"]
usage: "(cond [pred then]*)"
tags: ["conditional"]
---
For each *pred-then* clause, the predicate will be evaluated, and if it is truthy (not _#f_ (false) or the empty list) the *then* form is evaluated and returned, otherwise the next clause is processed.

#### An Example

~~~scheme
(define x 99)

(cond
  [(< x 50)  "was less than 50"]
  [(> x 100) "was greater than 100"]
  [:else     "was in between"])
~~~

In this case, _"was in between"_ will be returned. The reason that this works is because the `:else` keyword, like all keywords, evaluates to truthy.

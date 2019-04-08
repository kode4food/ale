---
title: "and"
date: 2019-04-06T12:19:22+02:00
description: "performs a short-circuiting and"
names: ["and"]
usage: "(and form*)"
tags: ["conditional", "logic"]
---
Evaluates the forms from left to right. As soon as one evaluates to a falsey value, will return that value. Otherwise it will proceed to evaluating the next form.

#### An Example

```clojure
(and (+ 1 2 3)
     false
     "not returned")
```

Will return _false_, never evaluating _"not returned"_, whereas:

```clojure
(and (+ 1 2 3)
     true
     "returned")
```

Will return the string _"returned"_.

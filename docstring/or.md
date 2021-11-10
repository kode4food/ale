---
title: "or"
description: "performs a short-circuiting or"
names: ["or"]
usage: "(or form*)"
tags: ["conditional", "logic"]
---

Evaluates the forms from left to right. As soon as one evaluates to a truthy value, will return that value. Otherwise, it will proceed to evaluating the next form.

#### An Example

```scheme
(or (+ 1 2 3)
    false
    "not returned")
```

Will return _6_, never evaluating `false` (false) and `"not returned"`, whereas:

```scheme
(or false
    '()
    "returned")
```

Will return the string _"returned"_.

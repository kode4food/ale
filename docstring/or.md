---
title: "or"
date: 2019-04-06T12:19:22+02:00
description: "performs a short-circuiting or"
names: ["or"]
usage: "(or form*)"
tags: ["conditional", "logic"]
---
Evaluates the forms from left to right. As soon as one evaluates to a truthy value, will return that value. Otherwise it will proceed to evaluating the next form.

#### An Example

~~~scheme
(or (+ 1 2 3)
    #f
    "not returned")
~~~

Will return _6_, never evaluating `#f` (false) and `"not returned"`, whereas:

~~~scheme
(or #f
    '()
    "returned")
~~~

Will return the string _"returned"_.

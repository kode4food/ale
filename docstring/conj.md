---
title: "conj"
date: 2019-04-06T12:19:22+02:00
description: "adds elements to a sequence"
names: ["conj"]
usage: "(conj seq form+)"
tags: ["sequence"]
---
Adds elements to a conjoinable sequence. This behavior will differ depending on the concrete type. A list will prepend, a vector will append, while an object makes no guarantees about ordering.

#### An Example

~~~scheme
(conj [1 2 3 4] 5 6 7 8)
~~~

Will return the vector _[1 2 3 4 5 6 7 8]_.

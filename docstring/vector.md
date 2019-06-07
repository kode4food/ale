---
title: "vector"
date: 2019-04-06T12:19:22+02:00
description: "creates a new vector"
names: ["vector"]
usage: "(vector form*)"
tags: ["sequence"]
---
Will create a new vector whose elements are the evaluated forms provided. This function is no different than the vector literal syntax except that it can be treated in a first-class fashion.

#### An Example

~~~scheme
(def x "hello")
(def y "there")
(vector x y)
~~~

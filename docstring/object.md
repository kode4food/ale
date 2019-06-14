---
title: "object"
date: 2019-04-06T12:19:22+02:00
description: "creates a new object (hash-map) instance"
names: ["object"]
usage: "(object <key value>*)"
tags: ["data"]
---
Will create a new hash map based on the provided key-value pairs, or return an empty object if no forms are provided. This function is no different than the object literal syntax except that it can be treated in a first-class fashion.

#### An Example

~~~scheme
(object :name "ale" :age 0.3 :lang "Go")

;; which is the same as
{:name "ale" :age 0.3 :lang "Go"}
~~~

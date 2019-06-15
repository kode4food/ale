---
title: "not"
date: 2019-04-06T12:19:22+02:00
description: "logically inverts the truthiness of the provided form"
names: ["not"]
usage: "(not form)"
tags: ["logic"]
---
Will return _#f_ (false) if the provided *form* is truthy, otherwise will return _#t_ (true).

#### An Example

~~~scheme
(not "hello")
~~~

This will return the boolean _#f_ (false) because the value _"hello"_ is truthy.

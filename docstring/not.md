---
title: "not"
date: 2019-04-06T12:19:22+02:00
description: "logically inverts the truthiness of the provided form"
names: ["not"]
usage: "(not form)"
tags: ["logic"]
---
Will return _false_ if the provided *form* is truthy, otherwise will return _true_.

#### An Example

~~~scheme
(not "hello")
~~~

This will return the boolean _false_ because the value _"hello"_ is truthy.

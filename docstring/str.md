---
title: "str"
description: "converts forms to a string"
date: 2019-04-06T12:19:22+02:00
names: ["str", "str!"]
usage: "(str form*) (str! form*)"
tags: ["sequence", "conversion"]
---
Creates a new string from the stringified values of the provided forms.

#### An Example

~~~scheme
(str "hello" [1 2 3 4])
~~~

This example will return the string _"hello[1 2 3 4]"_.

#### Reader Strings

Alternatively, one can use the `str!` function to produce a stringified version that *may* be able to be read by the Ale reader.

~~~scheme
(str! "hello" [1 2 3 4])
~~~

This example will return the string _"\"hello\" [1 2 3 4]"_.

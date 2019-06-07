---
title: "remainder (mod)"
date: 2019-04-06T12:19:22+02:00
description: "calculates the remainder of a number sequence"
names: ["mod"]
usage: "(mod form+)"
tags: ["math", "number"]
---
Takes a set of numbers and calculates the collective remainder of dividing each by the next.

#### An Example

~~~scheme
(mod 10 3)      ;; returns 1
(mod 20 7.0 4)  ;; returns 2
~~~

---
title: "fn"
date: 2019-04-06T12:19:22+02:00
description: "creates an anonymous function"
names: ["fn"]
usage: "(fn name? [args] form*)"
tags: ["function"]
---
Will create an anonymous function that may be passed around in a first-class manner.

#### An Example

~~~scheme
(define double
  (let [mul 2]
    (fn "doubles values" [x] (* x mul))))

(to-vector
  (map double '(1 2 3 4 5 6)))
~~~

This example will return the vector _[2 4 6 8 10 12]_.

Anonymous functions produce a closure that copies the bindings that have been referenced from the surrounding scope.

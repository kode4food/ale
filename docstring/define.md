---
title: "define"
date: 2019-04-06T12:19:22+02:00
description: "binds a namespace entry"
names: ["define"]
usage: "(define name form) (define (name param*) form*)"
tags: ["binding"]
---
Will bind a value to a global name. All bindings are immutable and result in an error being raised if an attempt is made to re-bind them. This behavior is different than most Lisps, as they will generally fail silently in such cases.

#### An Example

~~~scheme
(define x
  (map
    (lambda (y) (* y 2))
    seq1 seq2 seq3))
~~~

This example will create a lazy map where each element of the three provided sequences is doubled upon request. It will then bind that lazy map to the name *x*.

---
title: "lambda"
description: "creates a lambda"
names: ["lambda" "lambda-rec"]
usage: "(lambda (param*) form*) (lambda-rec name (param*) form*)"
tags: ["function"]
---

Will create a lambda function that may be passed around in a first-class manner.

#### An Example

```scheme
(define double
  (let [mul 2]
    (lambda (x) (* x mul))))

(seq->vector
  (map double '(1 2 3 4 5 6)))
```

This example will return the vector _[2 4 6 8 10 12]_.

Lambdas produce a closure that copies the bindings that have been referenced from the surrounding scope.

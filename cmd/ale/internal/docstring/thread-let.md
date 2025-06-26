---
title: "let thread (let->)"
description: "threaded let binding with intermediate value access"
names: ["let->"]
usage: "(let-> [binding-name initial-value] forms*)"
tags: ["function"]
---

A threaded let binding that threads the bound value through each form by name. Each form receives the current value of the binding and the result becomes the new value of the binding.

#### An Example

```scheme
(let-> [x 10]
       (+ x 5)
       (\* x 2)
       (/ x 2))
```

#### Another Example

```scheme
(let-> [user {}]
       (assoc user (:name . "John"))
       (assoc user (:age . 25))
       (assoc user (:status . "active")))
```

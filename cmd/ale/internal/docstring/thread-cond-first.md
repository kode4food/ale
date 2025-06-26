---
title: "conditional thread first (cond->)"
description: "-> macro with conditional forms"
names: ["cond->"]
usage: "(cond-> expr [test form]*)"
tags: ["function"]
---

Like `->`, but each form is paired with a test condition. The form is only applied if the test evaluates to true. Each test is evaluated with the current threaded value in scope.

#### An Example

```scheme
(cond-> 10
        [true   (+ 5)]   ; Always applies: 10 + 5 = 15
        [(> 12) (\* 2)]   ; Applies if > 12: 15 \* 2 = 30
        [(< 25) (/ 3)]   ; Doesn't apply (30 is not < 25)
        [true   (- 1)])  ; Always applies: 30 - 1 = 29
```

#### Another Example

```scheme
(cond-> {}
        [empty?            (assoc (:name . "John"))]
        [!empty?           (assoc (:size . "non-empty"))]
        [(contains? :name) (assoc (:greeting . "Hello"))]
        [object?           (assoc (:type . "object"))])
```

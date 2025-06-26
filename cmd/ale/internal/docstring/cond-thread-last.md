---
title: "conditional thread last (cond->>)"
description: "->> macro with conditional forms"
names: ["cond->>"]
usage: "(cond->> expr [test form]*)"
tags: ["function"]
---

Like `->>`, but each form is paired with a test condition. The form is only applied if the test evaluates to true. Each test is evaluated with the current threaded value in scope.

#### An Example

```scheme
(cond->> [1 2 3 4 5]
         [seq?    (map (lambda (x) (* x 2)))]  ; doubles: [2 4 6 8 10]
         [!empty? (filter even?)]              ; keeps evens: [2 4 6 8 10]
         [(> 2)   (take 3)]                    ; takes first 3: [2 4 6]
         [false   (concat [0])])               ; skipped
```

#### Another Example

```scheme
(cond->> [1 2 3]
         [seq? (filter even?)]    ; predicate true, filters to [2]
         [empty? (concat [4 5])]  ; predicate false, skips concat
         [true (take 1)])         ; takes first element
```

---
title: "car/cdr accessors"
description: "compound sequence accessors for nested pairs"
names: ["caar", "cadr", "cdar", "cddr", "caaar", "caadr", "cadar", "caddr", "cdaar", "cdadr", "cddar", "cdddr", "caaaar", "caaadr", "caadar", "caaddr", "cadaar", "cadadr", "caddar", "cadddr", "cdaaar", "cdaadr", "cdadar", "cdaddr", "cddaar", "cddadr", "cddar", "cdddar", "cddddr"]
usage: "(<accessor> seq)"
tags: ["function", "sequence"]
---

These are compound sequence accessors that combine multiple `car` and `cdr` operations. The name of each function describes which operations to perform from  right to left:

- `a` represents `car` (get first element)
- `d` represents `cdr` (get rest of list)

#### Examples

Given a nested list structure:

```scheme
(define x '((1 2) (3 4)))

(caar x)   ; gets car of (car x) -> 1
(cadr x)   ; gets car of (cdr x) -> (3 4)
(cdar x)   ; gets cdr of (car x) -> (2)
(cddr x)   ; gets cdr of (cdr x) -> ()

;; More deeply nested examples
(define y '((1 2) ((3 4) 5)))
(cadar y)  ; car of (cdr of (car y)) -> 2
(caadr y)  ; car of (car of (cdr y)) -> (3 4)
```

The functions can access up to 4 levels deep in a nested list structure. Reading from right to left helps understand the access pattern - for example, `caddr` means "get the first element (car) of the rest of the rest (cddr) of the list".

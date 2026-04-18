---
title: "macro"
description: "wraps a procedure as a macro"
names: ["macro"]
usage: "(macro func)"
tags: ["macro"]
---

Converts a procedure into a macro value. The procedure receives the unevaluated input forms and returns expanded code. In most user code, `define-macro` is the more convenient interface.

#### An Example

```scheme
(define twice
  (macro (lambda (form) `(+ ,form ,form))))
```

---
title: "macro expansion"
description: "expands macro calls without evaluating them"
names: ["macroexpand", "macroexpand-1"]
usage: "(macroexpand form) (macroexpand-1 form)"
tags: ["macro", "compiler"]
---

`macroexpand-1` expands a single macro layer. `macroexpand` expands repeatedly until the form is no longer a macro call.

#### An Example

```scheme
(macroexpand-1 '(when true 42))
```

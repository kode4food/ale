---
title: "namespaces"
description: "creates, imports, and inspects namespaces"
names: ["define-namespace", "import", "declared"]
usage: "(define-namespace name form*) (import ns spec*) (declared ns?)"
tags: ["namespace", "special"]
---

These forms work with Ale namespaces. `define-namespace` creates and populates a namespace. `import` brings public names from another namespace into scope, either all at once or through explicit names and aliases. `declared` returns the names declared in the current or specified namespace.

#### An Example

```scheme
(define-namespace demo
  (define answer 42))
(import demo answer)
```

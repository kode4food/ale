---
title: "asm"
description: "emits raw virtual machine instructions"
names: ["asm"]
usage: "(asm instruction*)"
tags: ["compiler", "special"]
---

Provides direct access to Ale's assembler syntax. This form is primarily used by the core library and low-level code that needs explicit control over VM instructions.

#### An Example

```scheme
(asm
  const 99)
```

This example emits code that returns the literal value `99`.

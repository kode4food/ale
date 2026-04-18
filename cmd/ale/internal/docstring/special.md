---
title: "special"
description: "builds a low-level special form"
names: ["special"]
usage: "(special param-cases body)"
tags: ["compiler", "special"]
---

Constructs a compiler special form directly from parameter cases and assembler instructions. This is the low-level mechanism used to define forms such as `if`, and is mainly useful for core language implementation work rather than regular user code.

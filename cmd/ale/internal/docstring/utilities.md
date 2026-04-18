---
title: "utility helpers"
description: "small convenience functions and macros"
names: ["identity", "constantly", "no-op", "thunk"]
usage: "(identity value) (constantly value) (no-op form*) (thunk form*)"
tags: ["function"]
---

These are small utility forms used throughout Ale programs. `identity` returns its argument unchanged. `constantly` returns a function that always yields the same value. `no-op` ignores its arguments and returns `null`. `thunk` wraps forms in a zero-argument function.

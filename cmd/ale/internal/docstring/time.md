---
title: "time"
description: "evaluates forms and prints their elapsed duration"
names: ["time"]
usage: "(time form*)"
tags: ["performance"]
---

Evaluates the provided forms, prints the elapsed duration to standard output, and returns the value of the final form. Durations under `1000` are printed in nanoseconds; longer durations are printed in milliseconds.

#### An Example

```scheme
(time (+ 1 2 3))
```

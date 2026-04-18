---
title: "with-open"
description: "opens values and guarantees cleanup"
names: ["with-open"]
usage: "(with-open [name value]* form*)"
tags: ["io", "macro"]
---

Binds one or more values and ensures their `:close` procedures are called after the body finishes, even if an error is raised. If a bound value does not expose a callable `:close`, cleanup for that binding becomes a no-op.

#### An Example

```scheme
(with-open [r (:some-resource factory)]
  (: r :read))
```

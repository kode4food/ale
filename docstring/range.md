---
title: "range"
date: 2019-04-06T12:19:22+02:00
description: "creates a range"
names: ["range"]
usage: "(range min max inc)"
tags: ["sequence"]
---

Creates a lazy sequence that presents the numbers from _min_ (inclusive) to _max_ (exclusive), by _increment_. All parameters are optional. _min_ defaults to _0_, _max_ defaults to _\*pos-inf\*_, and _step_ defaults to _1_. If only one argument is provided, it is treated as _max_.

#### An Example

```scheme
(seq->vector (take 5 (range 10 inf 5)))
```

Will return the vector _[10 15 20 25 30]_.

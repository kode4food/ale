---
title: "range"
date: 2019-04-06T12:19:22+02:00
description: "creates a range"
names: ["range"]
usage: "(range min max inc)"
tags: ["sequence"]
---
Creates a lazy sequence that presents the numbers from *min* (inclusive) to *max* (exclusive), by *increment*. All arguments are optional. *min* defaults to _0_, *max* defaults to _\*pos-inf\*_, and *step* defaults to _1_. If only one argument is provided, it is treated as *max*.

#### An Example

```clojure
(to-vector (take 5 (range 10 inf 5)))
```

Will return the vector _[10 15 20 25 30]_.

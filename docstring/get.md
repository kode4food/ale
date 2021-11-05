---
title: "get"
description: "retrieves a value by key"
names: ["get"]
usage: "(get seq key default?)"
tags: ["sequence"]
---

Returns the value within a sequence that is associated with the specified key. If the key does not exist within the sequence, then either the default value is returned, or an error is raised.

#### An Example

```scheme
(define robert {:name "Bob" :age 45})
(get robert :address "wrong")
```

This example returns _"wrong"_ because the object doesn't contain an :address property.

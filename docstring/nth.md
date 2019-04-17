---
title: "nth"
date: 2019-04-06T12:19:22+02:00
description: "retrieves a value by index"
names: ["nth"]
usage: "(nth form index default?)"
tags: ["sequence"]
---
Returns the value that can be found at the specified index of its sequence. If the index is out of the bounds of the sequence, then either the default value is returned or an error is raised. Keep in mind that indexes are zero-based.

#### An Example

```clojure
(def l '(1 2 3 4))
(nth l 4 "wrong")
```

This example returns _"wrong"_ because index 4 (the 5th index) is beyond the end of the specified list.

#### Indexed Sequence Application

Instead of using the `nth` function, indexed sequences such as lists and vectors can also have arguments applied directly to them.

```clojure
(def l '(1 2 3 4))
(l 4 "wrong")
```

This will yield the same result as the previous example.

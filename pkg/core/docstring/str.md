---
title: "str"
description: "converts forms to a string"

names: ["str", "str!"]
usage: "(str form*) (str! form*)"
tags: ["sequence", "conversion"]
---

Creates a new string from the stringified values of the provided forms.

#### An Example

```scheme
(str "hello" [1 2 3 4])
```

This example will return the string _"hello[1 2 3 4]"_.

#### Reader Strings

Alternatively, one can use the `str!` function to produce a stringified version that _may_ be able to be read by the Ale reader.

```scheme
(str! "hello" [1 2 3 4])
```

This example will return the string _"\"hello\" [1 2 3 4]"_.

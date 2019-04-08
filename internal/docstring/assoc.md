---
title: "assoc"
date: 2019-04-06T12:19:22+02:00
description: "creates a new associative structure"
names: ["assoc"]
usage: "(assoc <key value>*)"
tags: ["data"]
---
Will create a new associative based on the provided key-value pairs, or return an empty associative if no forms are provided. This function is no different than the associative literal syntax except that it can be treated in a first-class fashion.

#### An Example

```clojure
(assoc
  :name "ale"
  :age  0.3
  :lang "Go")
```
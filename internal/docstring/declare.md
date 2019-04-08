---
title: "declare"
date: 2019-04-06T12:19:22+02:00
description: "forward declares a binding"
names: ["declare"]
usage: "(declare <name>)"
tags: ["binding"]
---
Forward declare a binding. This means that the name will be known in the current namespace, but not yet assigned. This can be useful when two functions refer to one another.

#### An Example

```clojure
(declare is-odd-number)

(defn is-even-number [n]
  (cond (= n 0) true
        :else   (is-odd-number (- n 1))))

(defn is-odd-number [n]
  (cond (= n 0) false
        :else   (is-even-number (- n 1))))
```

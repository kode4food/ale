---
title: "try/catch"
date: 2019-04-06T12:19:22+02:00
description: "handles raised errors"
names: ["try", "catch", "finally"]
usage: "(try form* catch-clause* finally-clause?)"
---

Will evaluate the provided forms, immediately short-circuiting if an error is raised, at which point control is passed into the first `catch` clause whose `predicate` evaluates to _#t_ (true). If a `finally` clause is defined, control will be passed into it after a successful evaluation of the provided forms or after a `catch` clause is evaluated.

`catch-clause` is defined as `(catch [name predicate] form*)`

`finally-clause` is defined as `(finally form*)`

#### An Example

```scheme
(try
  (raise "hello!")
  (println "won't reach me")
  (catch [n inf?] (println "won't match me"))
  (catch [s str?] (println "was a string ->" s))
  (finally (println "done")))
```

This will print the following to the console.

```
was a string -> hello!
done
```

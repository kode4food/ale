---
title: "defmacro"
date: 2019-04-06T12:19:22+02:00
description: "binds a reader macro"
names: ["defmacro"]
usage: "(defmacro name [name*] form+)"
tags: ["function", "macro", "binding"]
---
Will bind a macro to a global name. A macro is expanded by the reader in order to alter the source code's data representation before it is evaluated.

#### An Example

~~~scheme
(defmacro cond clauses
  (when (seq clauses)
    (if (= 1 (length clauses))
      (clauses 0)
      (list 'ale/if
        (clauses 0) (clauses 1)
        (cons 'cond (rest (rest clauses)))))))
~~~

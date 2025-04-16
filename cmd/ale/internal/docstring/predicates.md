---
title: "predicates"
description: "type and value testing functions"
names: ["any?", "!any?", "appendable?", "!appendable?", "boolean?", "!boolean?", "cons?", "!cons?", "even?", "!even?", "inf?", "!inf?", "-inf?", "!-inf?", "local?", "!local?", "macro?", "!macro?", "nan?", "!nan?", "number?", "!number?", "odd?", "!odd?", "pair?", "!pair?", "procedure?", "!procedure?", "qualified?", "!qualified?", "resolved?", "!resolved?", "promise-forced?", "!promise-forced?", "reversible?", "!reversible?", "special?", "!special?", "symbol?", "!symbol?"]
usage: "(<predicate> value)"
tags: ["function", "predicate"]
---

These predicates test values for specific types or properties. Each predicate has a corresponding negated version prefixed with "!". For example, `number?` tests if a value is a number, while `!number?` tests if it is not a number.

#### Examples

```scheme
;; Type checking
(number? 42)      ; -> #t
(boolean? true)   ; -> #t
(symbol? 'abc)    ; -> #t
(procedure? +)    ; -> #t

;; Value properties
(nan? (/ 0.0 0)) ; -> #t
(inf? +inf)      ; -> #t
(-inf? -inf)     ; -> #t
(even? 2)        ; -> #t
(odd? 3)         ; -> #t

;; Sequences
(pair? '(1 . 2))         ; -> #t
(appendable? [1 2 3])    ; -> #t ; vectors are appendable
(reversible? [1 2 3])    ; -> #t ; vectors can be reversed
(!appendable? '(1 2 3))  ; -> #t ; lists are not appendable

;; Symbols
(local? 'x)        ; -> #t ; for non-qualified symbols
(qualified? 'ns/x) ; -> #t ; for namespace-qualified symbols

;; Special forms and macros
(special? lambda) ; -> #t  ; lambda is a special form
(macro? when)     ; -> #t  ; when is a macro

;; Promise state
(define p (delay (+ 1 2)))
(resolved? p)     ; -> #f  ; not yet computed
(force p)
(resolved? p)     ; -> #t  ; now computed
```

Note that each negated predicate (`!predicate?`) returns the opposite of its corresponding positive predicate. For example, `(!even? 3)` is equivalent to `(not (even? 3))`.

Some key predicate categories:

- Type checking: `number?`, `boolean?`, `symbol?`, `procedure?`
- Sequence properties: `appendable?`, `reversible?`
- Symbol properties: `local?`, `qualified?`
- Form types: `special?`, `macro?`
- Numeric properties: `even?`, `odd?`, `inf?`, `nan?`
- Promise state: `resolved?` (alias: `promise-forced?`)

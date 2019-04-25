;;;; ale bootstrap: builtin

(def-special do)
(def-special if)
(def-special let)
(def-special fn)

(def-builtin read)
(def-special eval)
(def-builtin is-eq)

;; globals

(def-special declare)
(def-special def)

;; basic predicates

(def-builtin is-nil)
(def-builtin is-atom)
(def-builtin is-keyword)

;; macros

(def-special quote)
(def-special defmacro)
(def-special macroexpand-1)
(def-special macroexpand)
(def-builtin is-macro)
(def-macro syntax-quote)

;; symbols

(def-builtin sym)
(def-builtin gensym)
(def-builtin is-symbol)
(def-builtin is-local)
(def-builtin is-qualified)

;; strings

(def-builtin str)
(def-builtin str!)
(def-builtin is-str)

;; sequences

(def-builtin seq)
(def-builtin first)
(def-builtin rest)
(def-builtin last)
(def-builtin cons)
(def-builtin conj)
(def-builtin len)
(def-builtin nth)
(def-builtin get)
(def-builtin assoc)
(def-builtin list)
(def-builtin vector)

(def-builtin is-seq)
(def-builtin is-empty)
(def-builtin is-len)
(def-builtin is-indexed)
(def-builtin is-assoc)
(def-builtin is-mapped)
(def-builtin is-list)
(def-builtin is-vector)

;; numeric

(def-builtin +)
(def-builtin -)
(def-builtin *)
(def-builtin /)
(def-builtin mod)

(def-builtin =)
(def-builtin !=)
(def-builtin >)
(def-builtin >=)
(def-builtin <)
(def-builtin <=)

(def-builtin is-pos-inf)
(def-builtin is-neg-inf)
(def-builtin is-nan)

;; functions

(def-builtin apply)
(def-builtin is-apply)
(def-builtin is-special)

;; concurrency

(def-builtin go*)
(def-builtin chan)
(def-builtin promise)
(def-builtin is-promise)

;; lazy sequences

(def-builtin lazy-seq*)
(def-builtin concat)
(def-builtin append)
(def-builtin filter)
(def-builtin map)
(def-builtin take)
(def-builtin drop)
(def-builtin reduce)
(def-builtin for-each*)

;; raise and recover

(def-builtin raise)
(def-builtin recover)
(def-builtin defer)

;; current time

(def-builtin current-time)

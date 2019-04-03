# (partial func arg+) generates a function based on a partial apply
Returns a new Function whose initial arguments are pre-bound to those provided. When that Function is invoked, any provided arguments will simply be appended to the pre-bound arguments before calling the original Function.

## An Example

  (def plus10 (partial + 4 6))
  (plus10 9)

This example will return _19_ as though `(+ 4 6 9)` were called.

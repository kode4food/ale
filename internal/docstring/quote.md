# (quote form) returns the specified form in data mode
Meaning that lists and symbols will not be evaluated. This macro is effectively the same as prepending an expression with an apostrophe (/'/).

## An Example

  (quote (1 2 3 4))

This will return the literal list rather than trying to apply the number 1, as if it were a function.  It is synonymous with the expression `'(1 2 3 4)`.

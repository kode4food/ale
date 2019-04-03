# (sym str) converts a string into an interned symbol
Converts the provided string into an interned symbol. Both qualified and local symbols are accepted.

## An Example

  (def hello-sym (sym "hello"))
  (eq hello-sym 'hello)

This example will return _true_.

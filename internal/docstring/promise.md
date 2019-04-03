# (promise) produces an unresolved value
Returns a promise in the form of a function. If applied without an argument, this function will *block*, waiting for a value to be delivered via a call to the promise function that includes an argument. A promise can only be delivered once.

## An Example

  (def p (promise))
  (p "hello")
  (p)

# (len seq) returns the length of the sequence
Will return the number of elements in a countable sequence. If the sequence is lazily computed, asynchronous, or otherwise incapable of being counted, this function will raise an error.

## An Example

  (len '(1 2 3 4))

This example will return _4_.

names: is-len len? !len?
# (len? form+) tests whether the provided forms are counted sequences
If all forms evaluate to a valid sequence than can report its length, then this function will return _true_. The first non-counted sequence will result in the function returning _false_.

## An Example

  (len? '(1 2 3 4) [5 6 7 8])

This example will return _true_.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not be valid counted sequences.

  (!len? "hello" 99)

This example will return _false_ because strings are counted sequences.

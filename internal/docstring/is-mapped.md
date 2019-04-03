names: is-mapped mapped? !mapped?
# (mapped? form+) tests whether the provided forms are mapped
If all forms evaluate to a mapped type, then this function will return _true_. The first non-mapped will result in the function returning _false_.

## An Example

  (mapped? {:name "bill"} {:name "peggy"} [1 2 3])

This example will return _false_ because the third form is a vector.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not be mapped.

  (!mapped? "hello" [1 2 3])

This example will return _true_.
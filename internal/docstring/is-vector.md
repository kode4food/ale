names: is-vector vector? !vector?
# (vector? form+) tests whether the provided forms are vectors
If all forms evaluate to a vector, then this function will return _true_. The first non-vector will result in the function returning _false_.

## An Example

  (vector? '(1 2 3 4) [5 6 7 8])

This example will return _false_ because the first form is a list.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not be vectors.

  (!vector? "hello" [99])

This example will return _false_ because the second form is a vector.

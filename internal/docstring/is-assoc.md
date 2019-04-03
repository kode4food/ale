names: is-assoc assoc? !assoc?
# (assoc? form+) tests whether the provided forms are associatives
If all forms evaluate to an assoc, then this function will return _true_. The first non-assoc will result in the function returning _false_.

## An Example

  (assoc? {:name "bill"} {:name "peggy"} [1 2 3])

This example will return _false_ because the third form is a vector.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not be associatives.

  (!assoc? "hello" [1 2 3])

This example will return _true_.

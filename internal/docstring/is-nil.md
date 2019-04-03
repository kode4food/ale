names: is-nil nil? !nil?
# (nil? form+) tests whether the provided forms are nil
If all forms evaluate to nil, then this function will return _true_. The first non-nil will result in the function returning _false_.

## An Example

  (nil? '(1 2 3 4) nil)

This example will return _false_ because the first form is a list.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not be nil.

  (!nil? "hello" [99])

This example will return _true_.

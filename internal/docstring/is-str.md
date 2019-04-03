names: is-str str? !str?
# (str? form+) tests whether the provided forms are strings
If all forms evaluate to strings, then this function will return _true_. The first non-string will result in the function returning _false_.

## An Example

  (str? '(1 2 3 4) "hello")

This example will return _false_ because the first form is a list.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not be strings.

  (!str? '(1 2 3) [99])

This example will return _true_.

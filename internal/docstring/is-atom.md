names: is-atom atom? !atom?
# (atom? form+) tests whether the provided forms are atomic
A form is considered to be atomic if it cannot be further evaluated and would otherwise evaluate to itself.

## An Example

  (atom? nil :hello "there")

This example will return _true_ because each value is atomic.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not be atomic.

  (!atom? '(+ 1 2 3) [4 5 6])

This example will return _true_ because compound types such as lists and vectors are not considered to be atomic.

names: eq !eq
# (eq form form+) tests if a set of values are identical to the first
Will return _false_ as soon as it encounters a form that is not identical to the first. Otherwise will return _true_.

## An Example

    (def h "hello")
    (eq "hello" h)

Like most predicates, this function can also be negated by prepending the /!/ character. In this case, _true_ will be returned if not all forms are equal.

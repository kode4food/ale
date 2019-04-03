# (when pred form**) conditionally evaluates a block
If the evaluated predicate is truthy (not _false_, not _nil_), the forms are evaluated. Will evaluate each form in turn, returning the final evaluation as its result.

## An Example

  (def x '(1 2 3 4 5 6 7 8))

  (when (> (len x) 3)
    (println "x is big")
    (len x))

If the symbol /when-not/ is used instead of /when/, then the predicate is evaluated and the block will be evaluated only if result is not truthy

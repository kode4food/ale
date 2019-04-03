# (fn name? [args] form+) creates an anonymous function
Will create an anonymous function that may be passed around in a first-class manner.

## An Example

  (def double
    (let [mul 2]
      (fn "doubles values" [x] (** x mul))))

  (to-vector
    (map double '(1 2 3 4 5 6)))

This example will return the vector _[2 4 6 8 10 12]_.

Anonymous functions produce a closure that copies the bindings that have been referenced from the surrounding scope.

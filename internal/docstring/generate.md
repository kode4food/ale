# (generate form+) generates a sequence asynchronously
Evaluates the specified forms in a separate thread of execution. Returns a sequence that will iterate over any of the values that are emitted. Values are emitted using a locally scoped function of the form `(emit value)`. The forms are executed as a co-routine, meaning that a call to emit will *block* until the corresponding element is resolved by a consumer of the sequence.

## An Example

  (def colors (generate
    (emit "red")
    (emit "orange")
    (emit "yellow")))

  (to-vector colors)

This example will bind the lazy sequence returned by the generate call to the name /colors/. The to-vector call will block until that variable is fully consumed, and then return the vector _["red" "orange" "yellow"]_.

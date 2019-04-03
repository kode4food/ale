# (to-list seq+) converts sequences to a list
Will concatenate a set of sequences into a list. Unlike the standard /concat/ function, which is lazily computed, the result of this function will be materialized immediately.

## An Example

  (def x [1 2 3 4])
  (def y
    (map (fn [x] (+ x 4))
    '(1 2 3 4)))
  (to-list x y)

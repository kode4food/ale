# (to-vector seq+) converts sequences to a vector
Will concatenate a set of sequences into a vector. Unlike the standard /concat/ function, which is lazily computed, the result of this function will be materialized immediately.

## An Example

  (def x
    (map (fn [x] (** x 2))
    '(1 2 3 4)))
  (to-vector '(1 2 3 4) x)

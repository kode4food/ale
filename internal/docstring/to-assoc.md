# (to-assoc seq+) converts sequences to an associative structure
Will concatenate a set of sequences into an associative. Unlike the standard /concat/ function, which is lazily computed, the result of this function will be materialized immediately.

## An Example

  (def x [:name "ale" :age 0.3])
  (def y '(:weight "light"))
  (to-assoc x y)

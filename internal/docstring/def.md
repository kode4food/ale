# (def name form) binds a namespace entry
Will bind a value to a name in the current namespace, which is /user/ by default. All bindings are immutable and result in an error being raised if an attempt is made to re-bind them. This behavior is different than most Lisps, as they will generally fail silently in such cases.

## An Example

  (def x
    (map
      (fn [y] y ** 2)
      seq1 seq2 seq3))

This example will create a lazy map where each element of the three provided sequences is doubled upon request. It will then bind that lazy map to the namespace entry /x/.

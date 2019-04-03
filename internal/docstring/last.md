# (last seq) returns the last element of the sequence
This function will return the last element of the specified sequence, or nil if the sequence is empty. It would be beneficial to check for a valid sequence using `(seq? seq)` before calling /last/.

## An Example

  (def x '(99 64 32 48))
  (last x)

 This example will return _48_.

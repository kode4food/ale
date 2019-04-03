# (first seq) returns the first element of the sequence
This function will return the first element of the specified sequence, or nil if the sequence is empty. It would be beneficial to check for a valid sequence using `(seq? seq)` before calling /first/ or /rest/.

## An Example

  (def x '(99 64 32 48))
  (first x)

 This example will return _99_.

# (rest seq) returns the rest of the sequence
This function will return a sequence that excludes the first element of the specified sequence. It would be beneficial to check for a valid sequence using `(seq? seq)` before calling /first/ or /rest/.

## An Example

  (def x '(99 64 32 48))
  (rest x)

 This example will return _(64 32 48)_.

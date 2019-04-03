# (go form**) asynchronously evaluates a block
The provided forms will be evaluated in a separate thread of execution. Any resulting value of the block will be discarded.

## An Example

  (def x (promise))
  (go (x "hello"))
  (str (x) " you!")

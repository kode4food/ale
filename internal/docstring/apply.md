# (apply func seq) applies arguments to a function
Evaluates the provided sequence and applies its values to the provided function (or applicable) as its arguments.

## An Example

  (def x '(1 2 3))
  (apply + x)

This example will return _6_.

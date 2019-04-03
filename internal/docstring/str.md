# (str form**) converts forms to a string
Creates a new string from the stringified values of the provided forms.

## An Example

  (str "hello" [1 2 3 4])

This example will return the string _hello[1 2 3 4]_.

## Reader Strings

Alternatively, one can use the `str!` function to produce a stringified version that *may* be able to be read by the Ale reader.

  (str! "hello" [1 2 3 4])

This example will return the string _"hello" [1 2 3 4]_.

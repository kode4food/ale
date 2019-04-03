# (get form key default?) retrieves a value by key
Returns the value within a sequence that is associated with the specified key. If the key does not exist within the sequence, then either the default value is returned or an error is raised.

## An Example

  (def robert {:name "Bob" :age 45})
  (get robert :address "wrong")

This example returns _"wrong"_ because the associative doesn't contain an :address property.

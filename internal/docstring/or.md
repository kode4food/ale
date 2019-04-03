# (or form**) performs a short-circuiting /or/
Evaluates the forms from left to right. As soon as one evaluates to a truthy value, will return that value. Otherwise it will proceed to evaluating the next form.

## An Example

  (or (+ 1 2 3)
      false
      "not returned")

Will return _6_, never evaluating _false_ and _"not returned"_, whereas:

  (or false
      nil
      "returned")

Will return the string _"returned"_.

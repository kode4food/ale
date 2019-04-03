# (and form**) performs a short-circuiting /and/
Evaluates the forms from left to right. As soon as one evaluates to a falsey value, will return that value. Otherwise it will proceed to evaluating the next form.

## An Example

  (and (+ 1 2 3)
       false
       "not returned")

Will return _false_, never evaluating _"not returned"_, whereas:

  (and (+ 1 2 3)
       true
       "returned")

Will return the string _"returned"_.

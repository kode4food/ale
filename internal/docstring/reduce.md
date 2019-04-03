# (reduce func val? seq) reduces a sequence
Iterates over a set of sequence, reducing their elements to a single resulting value. The function provided must take two arguments. The first and second sequence elements encountered are the initial values applied to that function. Thereafter, the result of the previous calculation is used as the first argument, while the next element is used as the second argument.

## An Example

  (reduce + 5 (range 1 11))

This will return the value _60_.

# (conj seq form+) adds elements to a sequence
Adds elements to a conjoinable sequence. This behavior will differ depending on the concrete type. A list will prepend, a vector will append, while an associative makes no guarantees about ordering. This function will not work with lazy sequences such as ones produced by `map` or `filter`.

## An Example

  (conj [1 2 3 4] 5 6 7 8)

Will return the vector _[1 2 3 4 5 6 7 8]_.

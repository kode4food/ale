# (concat seq+) lazily concatenates sequences
Creates a lazy sequence whose content is the result of concatenating the elements of each provided sequence.

## An Example

  (to-list (concat [1 2 3] '(4 5 6)))

This will return the list _(1 2 3 4 5 6)_

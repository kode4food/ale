# (filter func seq) lazily filters a sequence
Creates a lazy sequence whose content is the result of applying the provided function to the elements of the provided sequence. If the result of the application is truthy (not _false_, not _nil_) then the value will be included in the resulting sequence.

## An Example

  (filter (fn [x] (< x 3)) [1 2 3 4])

This will return the lazy sequence _(1 2)_

# (seq form) attempts to convert form to a sequence
Will attempt to convert the provided form to a sequence if it isn't already. If the form cannot be converted, _nil_ will be returned.

## An Example

  (when-let [s (seq "hello")]
    (to-vector (map (fn [x] (str x "-")) s)))

This example will return _["h-" "e-" "l-" "l-" "o-"]_.

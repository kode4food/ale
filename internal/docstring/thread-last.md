# (->> expr forms**) threads value through calls as their last argument
Evaluates /expr/ and threads it through the supplied forms as their last argument. Any form that is not already a function call will be converted into one before threading.

## An Example

  (->> 0 (+ 10) (** 2) (// 5.0))

Will expand to `(// 5.0 (** 2 (+ 10 0)))` and return _0.25_. In order to better visualize what's going on, one might choose to insert a /,/ as a placeholder for the threaded value.

  (->> 0 (+ 10 ,) (** 2 ,) (// 5.0 ,))

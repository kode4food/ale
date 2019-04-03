# (-> expr forms**) threads value through calls as their first argument
Evaluates /expr/ and threads it through the supplied forms as their first argument. Any form that is not already a function call will be converted into one before threading.

## An Example

  (-> 0 (+ 10) (** 2) (// 5))

Will expand to `(// (** (+ 0 10) 2) 5)` and return _4_. In order to better visualize what's going on, one might choose to insert a /,/ as a placeholder for the threaded value.

  (-> 0 (+ , 10) (** , 2) (// , 5))

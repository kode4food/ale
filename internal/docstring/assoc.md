# (assoc <key value>**) creates a new associative structure
Will create a new associative based on the provided key-value pairs, or return an empty associative if no forms are provided. This function is no different than the associative literal syntax except that it can be treated in a first-class fashion.

## An Example

  (assoc
    :name "ale"
    :age  0.3
    :lang "Go")

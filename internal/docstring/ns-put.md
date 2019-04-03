# (ns-put ns name form) binds a namespace entry
Will bind a value to a name in the specified namespace. All bindings are immutable and result in an error being raised if an attempt is made to re-bind them.

## An Example

  (ns-put (ns config) env "production")

This example will bind the string _"production"_ to the name /env/ in the /config/ namespace.

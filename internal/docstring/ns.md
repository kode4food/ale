# (ns name?) returns a namespace
Will return a namespace environment based on its name. If the namespace does not already exist, it will be created before being returned. If no name is provided, the current namespace will be returned.

## An Example

  (let [n (ns some-new-namespace)]
    (prn n))

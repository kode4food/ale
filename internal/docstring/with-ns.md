# (with-ns name form+) evaluates forms against a namespace
Evaluates a set of forms in the environment of a specified namespace. If the namespace does not already exist, it will be created before evaluating the forms. The evaluated result of the last form will be returned.

## An Example

  (def x "outside the namespace")
  (with-ns my-namespace
    (def x "inside the namespace")
    (prn "print 1: " x))
  (prn "print 2:" x)

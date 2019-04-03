# (partition count step? seq) partitions a sequence
Will partition a sequence into groups of /count/ elements, incrementing by the number of elements defined in /step/ (or /count/ if /step/ is not provided).

## An Example

  (to-list (partition 2 3 [1 2 3 4 5 6 7 8 9 10]))

This example will return _((1 2) (4 5) (7 8) (10))_.

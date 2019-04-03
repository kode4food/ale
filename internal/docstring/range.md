# (range min max inc) creates a range
Creates a lazy sequence that presents the numbers from /min/ (inclusive) to /max/ (exclusive), by /increment/. All arguments are optional. /min/ defaults to _0_, /max/ defaults to _inf_, and /step/ defaults to _1_. If only one argument is provided, it is treated as /max/.

## An Example

    (to-vector (take 5 (range 10 inf 5)))

Will return the vector _[10 15 20 25 30]_.

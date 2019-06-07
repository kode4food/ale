---
title: "partition"
date: 2019-04-06T12:19:22+02:00
description: "partitions a sequence"
names: ["partition"]
usage: "(partition count step? seq)"
tags: ["sequence", "comprehension"]
---
Will partition a sequence into groups of *count* elements, incrementing by the number of elements defined in *step* (or *count* if *step* is not provided).

#### An Example

~~~scheme
(to-list (partition 2 3 [1 2 3 4 5 6 7 8 9 10]))
~~~

This example will return _((1 2) (4 5) (7 8) (10))_.

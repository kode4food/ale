
---
title: "chan"
date: 2019-04-06T12:19:22+02:00
description: "creates a unidirectional channel"
names: ["chan"]
usage: "(chan size?)"
tags: ["concurrency"]
---

A channel is a data structure that is used to generate a lazy sequence of values. The result is a hash-map consisting of an `emit` function, a `close` function, and a sequence. Depending on the size of the channel's buffer, retrieving an element from the sequence _may block_, waiting for the next value to be emitted or for the channel to be closed. Emitting a value to a channel will also block until the buffer is flushed as a result of iterating over the sequence.

#### Channel Keys

```
*:seq*     the sequence to be generated
*:emit*    an emitter function of the form (emit value+)
*:close*   a function to close the channel (close)
```

#### An Example

```scheme
(let [ch (chan)]
  (go (: ch :emit "foo")
      (: ch :emit "bar")
      (: ch :close))

  (seq->vector (ch :seq)))
```

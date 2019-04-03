# (chan) creates a unidirectional channel
A channel is a data structure that is used to generate a lazy sequence of values. The result is a hash-map consisting of an `emit` function, a `close` function, and a sequence. Retrieving an element from the sequence may *block*, waiting for the next value to be emitted or for the channel to be closed. Emitting a value to a channel will also block until the buffer is flushed as a result of iterating over the sequence.

## Channel Keys

*:seq*     the sequence to be generated
*:emit*    an emitter function of the form `(emit value)`
*:close*   a function to close the channel `(close)`

## An Example

  (let [ch (chan)]
    (go (. ch :emit "foo")
        (. ch :emit "bar")
        (. ch :close))

    (to-vector (:seq ch)))

Channels represent the low-level underpinnings of the `generate` function, which is in most cases the preferred way to generate lazy sequences. The example above could be rewritten:

  (to-vector (generate
    (emit "foo")
    (emit "bar")))

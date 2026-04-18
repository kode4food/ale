---
title: "spawn"
description: "starts an actor with a mailbox"
names: ["spawn"]
usage: "(spawn func mbox-size? monitor?)"
tags: ["concurrency"]
---

Starts an actor-like process whose function receives a mailbox sequence. The returned value is a sender procedure that can emit messages into that mailbox. The default mailbox size is `16`. If a `monitor` procedure is provided, it is used to handle raised values from the actor body.

#### An Example

```scheme
(define send
  (spawn (lambda (mailbox) (first mailbox))))
```

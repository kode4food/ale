---
title: "seq!"
description: "asserts that a value can act as a non-empty sequence"
names: ["seq!"]
usage: "(seq! value)"
tags: ["sequence"]
---

Returns `value` when it can act as a non-empty sequence. If `value` is a sequence but empty, returns _#f_. If `value` cannot act as a sequence at all, an error is raised.

---
title: "sequence combinators"
description: "derived sequence transformation helpers"
names: ["zip", "mapcat", "cartesian-product", "map!", "take-while"]
usage: "(zip seq+) (mapcat func seq+) (cartesian-product seq+) (map! func seq) (take-while pred seq)"
tags: ["sequence"]
---

These helpers build or transform sequences. `zip` groups aligned values into lists and stops when the shortest input ends. `mapcat` maps, then lazily concatenates the mapped results. `cartesian-product` produces every combination across the provided sequences. `map!` eagerly maps into a list. `take-while` lazily consumes values while a predicate stays true.

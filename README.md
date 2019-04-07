# Ale is a Lisp Environment
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/kode4food/ale)](https://goreportcard.com/report/gitlab.com/kode4food/ale) [![Build Status](https://travis-ci.org/kode4food/ale.svg?branch=master)](https://travis-ci.org/kode4food/ale) [![Coverage Status](https://coveralls.io/repos/gitlab/kode4food/ale/badge.svg)](https://coveralls.io/gitlab/kode4food/ale)

Ale is a Lisp Environment written in [Go](https://golang.org/).

## How To Install
Make sure your `GOPATH` is set, then run `go get` to retrieve the package.

```bash
go get gitlab.com/kode4food/ale/cmd/ale
```

## How To Run A Source File
Once you've installed the package, you can run it from `GOPATH/bin` like so:

```bash
ale somefile.lisp

# or

cat somefile.lisp | ale
```

## How To Start The REPL
Ale has a very crude Read-Eval-Print Loop that will be more than happy
to start if you invoke `ale` with no arguments from your shell.

## Current Status
Still a work in progress, and the compiler is pretty fragile, but that will
change rapidly. Fine out more at [Ale's Home](https://www.ale-lang.org)

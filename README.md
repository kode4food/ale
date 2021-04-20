# Ale is a Lisp Environment

[![Go Report Card](https://goreportcard.com/badge/github.com/kode4food/ale)](https://goreportcard.com/report/github.com/kode4food/ale) [![Build Status](https://travis-ci.org/kode4food/ale.svg?branch=master)](https://travis-ci.org/kode4food/ale) [![Test Coverage](https://api.codeclimate.com/v1/badges/bcf86d6aa52ebaaed63f/test_coverage)](https://codeclimate.com/github/kode4food/ale/test_coverage) [![Maintainability](https://api.codeclimate.com/v1/badges/bcf86d6aa52ebaaed63f/maintainability)](https://codeclimate.com/github/kode4food/ale/maintainability) [![GitHub](https://img.shields.io/github/license/kode4food/ale)](https://github.com/kode4food/ale/blob/master/LICENSE.md)

Ale is a Lisp Environment for [Go](https://golang.org/) applications

## How To Install

Make sure your `GOPATH` is set, then run `go get` to retrieve the package.

```bash
go get github.com/kode4food/ale/cmd/ale
```

## How To Run A Source File

Once you've installed the package, you can run it from `GOPATH/bin` like so:

```bash
ale somefile.ale

# or

cat somefile.ale | ale
```

## How To Start The REPL

Ale has a very crude Read-Eval-Print Loop that will be more than happy
to start if you invoke `ale` with no arguments from your shell.

## Current Status

Still a work in progress, and the compiler is pretty fragile, but that will
change rapidly.

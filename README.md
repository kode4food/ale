# Ale is a Lisp Environment

[![Go Report Card](https://goreportcard.com/badge/github.com/kode4food/ale)](https://goreportcard.com/report/github.com/kode4food/ale) [![Build Status](https://app.travis-ci.com/kode4food/ale.svg?branch=main)](https://app.travis-ci.com/kode4food/ale) [![Test Coverage](https://api.codeclimate.com/v1/badges/bcf86d6aa52ebaaed63f/test_coverage)](https://codeclimate.com/github/kode4food/ale/test_coverage) [![Maintainability](https://api.codeclimate.com/v1/badges/bcf86d6aa52ebaaed63f/maintainability)](https://codeclimate.com/github/kode4food/ale/maintainability) [![GitHub](https://img.shields.io/github/license/kode4food/ale)](https://github.com/kode4food/ale/blob/main/LICENSE.md)

Ale is a Lisp Environment for [Go](https://golang.org/) applications

## How To Install

Make sure your `GOPATH` is set, then run `go install` to install the command line tool.

```bash
go install github.com/kode4food/ale/cmd/ale@latest
```

## How To Run A Source File

Once you've installed the package, you can run it from `GOPATH/bin` like so:

```bash
ale somefile.ale

# or

cat somefile.ale | ale


## or

ale <<EOF
(let [ch (chan)]
  (go (: ch :emit "Hello")
      (: ch :emit ", ")
      (: ch :emit "Ale!")
      (: ch :close))

  (println (apply str (ch :seq))))
EOF
```

## How To Start The REPL

Ale has a very crude Read-Eval-Print Loop that will be more than happy
to start if you invoke `ale` with no arguments from your shell.

## Current Status

Still a work in progress. Use at your own risk.

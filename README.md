# Ale is a Lisp Environment
[![Go Report Card](https://goreportcard.com/badge/github.com/kode4food/ale)](https://goreportcard.com/report/github.com/kode4food/ale) [![Build Status](https://travis-ci.org/kode4food/ale.svg?branch=master)](https://travis-ci.org/kode4food/ale) [![Coverage Status](https://coveralls.io/repos/github/kode4food/ale/badge.svg?branch=master)](https://coveralls.io/github/kode4food/ale?branch=master) [![Maintainability](https://api.codeclimate.com/v1/badges/f1eff5eeb0ae12973b4a/maintainability)](https://codeclimate.com/github/kode4food/ale/maintainability)  [![Sponsor](https://img.shields.io/badge/❤️%20Sponsor%20-%20Patreon-f39f37)](https://www.patreon.com/ale_lang)

Ale is a Lisp Environment for [Go](https://golang.org/) applications

## How To Install
Make sure your `GOPATH` is set, then run `go get` to retrieve the package.

```bash
go get github.com/kode4food/ale/cmd/ale
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
change rapidly.

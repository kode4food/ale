package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/read"
)

// Error messages
const (
	ErrFileNotFound = "file not found: %s"
)

// EvaluateStdIn reads from StdIn and evaluates it
func EvaluateStdIn() {
	defer exitWithError()

	buffer, _ := ioutil.ReadAll(os.Stdin)
	evalBuffer(ns, buffer)
}

// EvaluateFile reads the specific source file and evaluates it
func EvaluateFile() {
	defer exitWithError()

	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err != nil {
		fmt.Println(fmt.Errorf(ErrFileNotFound, filename))
		os.Exit(-1)
	} else {
		evalBuffer(ns, buffer)
	}
}

func evalBuffer(ns env.Namespace, src []byte) data.Value {
	r := read.FromString(data.String(src))
	return eval.Block(ns, r)
}

func exitWithError() {
	if rec := recover(); rec != nil {
		if ev, ok := rec.(error); ok {
			fmt.Println(ev.Error())
		} else {
			fmt.Println(rec)
		}
		os.Exit(-2)
	}
}

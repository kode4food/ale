package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/read"
)

const fileNotFound = "file not found: %s"

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
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		evalBuffer(ns, buffer)
	} else {
		fmt.Println(fmt.Sprintf(fileNotFound, filename))
		os.Exit(-1)
	}
}

func evalBuffer(ns namespace.Type, src []byte) api.Value {
	r := read.FromString(api.String(src))
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

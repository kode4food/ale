package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/read"
)

const (
	//ErrFileNotFound is raised when a file is not found
	ErrFileNotFound = "file not found: %s"
)

// EvaluateStdIn reads from StdIn and evaluates it
func EvaluateStdIn() {
	defer exitWithError()

	buffer, _ := io.ReadAll(os.Stdin)
	evalBuffer(buffer)
}

// EvaluateFile reads the specific source file and evaluates it
func EvaluateFile(filename string) {
	defer exitWithError()

	buffer, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(fmt.Errorf(ErrFileNotFound, filename))
		os.Exit(-1)
	}
	evalBuffer(buffer)
}

func evalBuffer(src []byte) data.Value {
	ns := makeUserNamespace()
	r := read.FromString(data.String(src))
	return eval.Block(ns, r)
}

func exitWithError() {
	if rec := recover(); rec != nil {
		if ev, ok := rec.(error); ok {
			fmt.Println(ev.Error())
		}
		os.Exit(-2)
	}
}

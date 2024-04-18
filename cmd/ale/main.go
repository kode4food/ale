package main

import (
	"os"

	"github.com/kode4food/ale/cmd/ale/internal"
)

func main() {
	switch {
	case isStdInPiped():
		internal.EvaluateStdIn()
	case len(os.Args) < 2:
		internal.NewREPL().Run()
	default:
		internal.EvaluateFile(os.Args[1])
	}
}

func isStdInPiped() bool {
	s, _ := os.Stdin.Stat()
	return (s.Mode() & os.ModeCharDevice) == 0
}

package main

import "os"

func main() {
	switch {
	case isStdInPiped():
		EvaluateStdIn()
	case len(os.Args) < 2:
		NewREPL().Run()
	default:
		EvaluateFile(os.Args[1])
	}
}

func isStdInPiped() bool {
	s, _ := os.Stdin.Stat()
	return (s.Mode() & os.ModeCharDevice) == 0
}

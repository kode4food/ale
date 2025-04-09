package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/console"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/env"
)

// REPL manages a Read-Eval-Print Loop
type REPL struct {
	ns  env.Namespace
	rl  *readline.Instance
	buf strings.Builder
	idx int
}

// NewREPL instantiates a new REPL instance
func NewREPL() *REPL {
	repl := &REPL{
		ns: makeUserNamespace(),
	}
	repl.registerBuiltIns()

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile:            getHistoryFile(),
		Painter:                console.Painter(),
		DisableAutoSaveHistory: true,
		AutoComplete:           repl,
	})
	if err != nil {
		panic(debug.ProgrammerError(err.Error()))
	}

	repl.rl = rl
	repl.idx = 1

	return repl
}

// Run will perform the Read-Eval-Print-Loop
func (r *REPL) Run() {
	defer func() {
		_ = r.rl.Close()
	}()

	fmt.Println(ale.AppName, ale.Version)
	help()
	r.setInitialPrompt()

	for {
		line, err := r.rl.Readline()
		r.buf.WriteString(line + "\n")
		fmt.Print(console.Reset)

		if err != nil {
			emptyBuffer := isEmptyString(r.buf.String())
			if errors.Is(err, readline.ErrInterrupt) && !emptyBuffer {
				r.reset()
				continue
			}
			break
		}

		if isEmptyString(line) {
			continue
		}

		if !r.evalBuffer() {
			r.setContinuePrompt()
			continue
		}

		r.saveHistory()
		r.reset()
	}
	shutdown()
}

func (r *REPL) reset() {
	r.buf.Reset()
	r.idx++
	r.setInitialPrompt()
}

// GetNS allows the tests to get at the namespace
func (r *REPL) GetNS() env.Namespace {
	return r.ns
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

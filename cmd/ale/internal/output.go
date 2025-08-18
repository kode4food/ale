package internal

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/cmd/ale/internal/console"
)

type sentinel struct{}

const (
	domain = console.Domain + "%s" + console.Reset + " "
	prompt = domain + "[%d]> " + console.Code
	cont   = domain + "[%d]" + console.NewLine + "   " + console.Code

	output = console.Bold + "%s" + console.Reset
	good   = domain + console.Result + "[%d]= " + output
	bad    = domain + console.Error + "[%d]! " + output

	stripped = "%s [%d]= "
)

var nothing = sentinel{}

func (r *REPL) setInitialPrompt() {
	name := r.ns.Domain()
	r.setPrompt(fmt.Sprintf(prompt, name, r.idx))
}

func (r *REPL) setContinuePrompt() {
	r.setPrompt(fmt.Sprintf(cont, r.nsSpace(), r.idx))
}

func (r *REPL) setPrompt(s string) {
	r.rl.SetPrompt(s)
}

func (r *REPL) nsSpace() string {
	ns := string(r.ns.Domain())
	return strings.Repeat(" ", len(ns))
}

func (r *REPL) outputResult(v ale.Value) {
	if v == nothing {
		return
	}

	domain := string(r.ns.Domain())
	offset := len(fmt.Sprintf(stripped, domain, r.idx))

	sv := PrettyPrintAt(v, offset)
	res := fmt.Sprintf(good, r.nsSpace(), r.idx, sv)
	fmt.Println(res)
}

func (r *REPL) outputError(err error) {
	msg := err.Error()
	res := fmt.Sprintf(bad, r.nsSpace(), r.idx, msg)
	fmt.Println(res)
}

func (sentinel) Equal(ale.Value) bool {
	return false
}

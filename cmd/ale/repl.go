package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/docstring"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/console"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/markdown"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/read"
	"github.com/kode4food/ale/read/lex"
	"github.com/kode4food/ale/read/parse"
	"github.com/kode4food/comb/slices"
)

type (
	sentinel struct{}

	// REPL manages a FromLexer-Eval-Print Loop
	REPL struct {
		buf bytes.Buffer
		rl  *readline.Instance
		ns  env.Namespace
		idx int
	}
)

const (
	// UserDomain is the name of the namespace that the REPL starts in
	UserDomain = data.Local("user")

	domain = console.Domain + "%s" + console.Reset + " "
	prompt = domain + "[%d]> " + console.Code
	cont   = domain + "[%d]" + console.NewLine + "   " + console.Code

	output = console.Bold + "%s" + console.Reset
	good   = domain + console.Result + "[%d]= " + output
	bad    = domain + console.Error + "[%d]! " + output
)

var (
	anyChar     = regexp.MustCompile(lang.AnyChar)
	notPaired   = fmt.Sprintf(parse.ErrPrefixedNotPaired, "")
	docTemplate = docstring.MustGet("doc")
	nothing     = new(sentinel)
)

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
		panic(err)
	}

	repl.rl = rl
	repl.idx = 1

	return repl
}

// Run will perform the Eval-Print-Loop
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
			if err == readline.ErrInterrupt && !emptyBuffer {
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
	return anyChar.ReplaceAllString(ns, " ")
}

func (r *REPL) evalBuffer() (completed bool) {
	defer func() {
		if err := toError(recover()); err != nil {
			if isRecoverable(err) {
				completed = false
				return
			}
			r.outputError(err)
			completed = true
		}
	}()

	seq := r.scanBuffer()
	r.evalBlock(r.ns, seq)
	return true
}

func (r *REPL) scanBuffer() data.Sequence {
	src := data.String(r.buf.String())
	return read.FromString(src)
}

func (r *REPL) evalBlock(ns env.Namespace, seq data.Sequence) {
	v := sequence.ToVector(seq) // will trigger early reader error
	for _, f := range v {
		r.evalForm(ns, f)
	}
}

func (r *REPL) evalForm(ns env.Namespace, f data.Value) {
	defer func() {
		if err := toError(recover()); err != nil {
			r.outputError(err)
		}
	}()

	res := eval.Value(ns, f)
	r.outputResult(res)
}

func (r *REPL) outputResult(v data.Value) {
	if v == nothing {
		return
	}
	sv := data.MaybeQuoteString(v)
	res := fmt.Sprintf(good, r.nsSpace(), r.idx, sv)
	fmt.Println(res)
}

func (r *REPL) outputError(err error) {
	msg := err.Error()
	res := fmt.Sprintf(bad, r.nsSpace(), r.idx, msg)
	fmt.Println(res)
}

func (r *REPL) saveHistory() {
	defer func() { _ = recover() }()
	seq := r.scanBuffer()
	hist := string(sequence.ToStr(seq))
	_ = r.rl.SaveHistory(hist)
}

func (*sentinel) Equal(data.Value) bool {
	return false
}

func (*sentinel) String() string {
	return ""
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func toError(i any) error {
	if i == nil {
		return nil
	}
	switch i := i.(type) {
	case error:
		return i
	case data.Value:
		return errors.New(i.String())
	default:
		// Programmer error
		panic(fmt.Sprintf("non-standard error: %s", i))
	}
}

func isRecoverable(err error) bool {
	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		err = e
	}
	msg := err.Error()
	return msg == parse.ErrListNotClosed ||
		msg == parse.ErrVectorNotClosed ||
		msg == parse.ErrMapNotClosed ||
		msg == lex.ErrStringNotTerminated ||
		strings.HasPrefix(msg, notPaired)
}

// GetNS allows the tests to get at the namespace
func (r *REPL) GetNS() env.Namespace {
	return r.ns
}

func (r *REPL) getBuiltInsNamespace() env.Namespace {
	return r.ns.Environment().GetRoot()
}

func (r *REPL) registerBuiltIns() {
	r.registerBuiltIn("cls", data.MakeLambda(cls, 0))
	r.registerBuiltIn("doc", doc)
	r.registerBuiltIn("debug", data.MakeLambda(debugInfo, 0))
	r.registerBuiltIn("help", data.MakeLambda(help, 0))
	r.registerBuiltIn("quit", data.MakeLambda(shutdown, 0))
	r.registerBuiltIn("use", r.makeUse())
}

func (r *REPL) registerBuiltIn(n data.Local, v data.Value) {
	ns := r.getBuiltInsNamespace()
	ns.Declare(n).Bind(v)
}

func (r *REPL) makeUse() data.Value {
	return special.Call(func(e encoder.Encoder, args ...data.Value) {
		data.AssertFixed(1, len(args))
		n := args[0].(data.Local)
		old := r.ns
		r.ns = r.ns.Environment().GetQualified(n)
		if old != r.ns {
			fmt.Println()
		}
		generate.Literal(e, nothing)
	})
}

func shutdown(...data.Value) data.Value {
	t := time.Now().UTC().UnixNano()
	rs := rand.NewSource(t)
	rg := rand.New(rs)
	idx := rg.Intn(len(farewells))
	fmt.Println(farewells[idx])
	os.Exit(0)
	return nothing
}

func debugInfo(...data.Value) data.Value {
	runtime.GC()
	fmt.Println("Number of goroutines: ", runtime.NumGoroutine())
	return nothing
}

func cls(...data.Value) data.Value {
	fmt.Println(console.Clear)
	return nothing
}

func formatForREPL(s string) string {
	md := markdown.FormatMarkdown(s)
	lines := strings.Split(md, "\n")
	var out []string
	out = append(out, "")
	for _, l := range lines {
		if isEmptyString(l) {
			out = append(out, l)
		} else {
			out = append(out, "  "+l)
		}
	}
	out = append(out, "")
	return strings.Join(out, "\n")
}

func help(...data.Value) data.Value {
	md, err := docstring.Get("help")
	if err != nil {
		// Programmer error, help is missing
		panic(err)
	}
	fmt.Println(formatForREPL(md))
	return nothing
}

var doc = special.Call(func(e encoder.Encoder, args ...data.Value) {
	data.AssertRanged(0, 1, len(args))
	if len(args) == 0 {
		docSymbolList()
	} else {
		docSymbol(args[0].(data.Local))
	}
	generate.Literal(e, nothing)
})

func docSymbol(sym data.Symbol) {
	name := string(sym.Name())
	if name == "doc" {
		docSymbolList()
		return
	}
	docStr := docstring.MustGet(name)
	f := formatForREPL(docStr)
	fmt.Println(f)
}

func docSymbolList() {
	names := docstring.Names()
	names = escapeNames(names)
	joined := strings.Join(names, ", ")
	f := formatForREPL(fmt.Sprintf(docTemplate, joined))
	fmt.Println(f)
}

var escapeNames = slices.Map(func(n string) string {
	if strings.Contains("`*_", n[0:1]) {
		return "\\" + n
	}
	return n
}).Must()

func makeUserNamespace() env.Namespace {
	return bootstrap.TopLevelEnvironment().GetQualified(UserDomain)
}

func getHistoryFile() string {
	if usr, err := user.Current(); err == nil {
		return path.Join(usr.HomeDir, ".ale-history")
	}
	return ""
}

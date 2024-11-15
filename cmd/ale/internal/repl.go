package internal

import (
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
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/internal/console"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/markdown"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/core/docstring"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/read"
	"github.com/kode4food/ale/pkg/read/lex"
	"github.com/kode4food/ale/pkg/read/parse"
	"github.com/kode4food/comb/slices"
)

type (
	sentinel struct{}

	// REPL manages a FromLexer-MustEval-Print Loop
	REPL struct {
		ns  env.Namespace
		rl  *readline.Instance
		buf strings.Builder
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

	escapeNames = slices.Map(func(n string) string {
		if strings.Contains("`*_", n[0:1]) {
			return "\\" + n
		}
		return n
	}).Must()
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
		panic(debug.ProgrammerError(err.Error()))
	}

	repl.rl = rl
	repl.idx = 1

	return repl
}

// Run will perform the MustEval-Print-Loop
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
	handleError := func(err error) bool {
		if isRecoverable(err) {
			return false
		}
		r.outputError(err)
		return true
	}

	defer func() {
		if err := toError(recover()); err != nil {
			completed = handleError(err)
		}
	}()

	seq := r.scanBuffer()
	if err := r.evalBlock(r.ns, seq); err != nil {
		return handleError(err)
	}
	return true
}

func (r *REPL) scanBuffer() data.Sequence {
	src := data.String(r.buf.String())
	return read.FromString(src)
}

func (r *REPL) evalBlock(ns env.Namespace, seq data.Sequence) error {
	v := sequence.ToVector(seq) // will trigger early reader error
	for _, f := range v {
		if err := r.evalForm(ns, f); err != nil {
			return err
		}
	}
	return nil
}

func (r *REPL) evalForm(ns env.Namespace, f data.Value) error {
	defer func() {
		if err := toError(recover()); err != nil {
			r.outputError(err)
		}
	}()

	res, err := eval.Value(ns, f)
	if err != nil {
		return err
	}
	r.outputResult(res)
	return nil
}

func (r *REPL) outputResult(v data.Value) {
	if v == nothing {
		return
	}
	sv := data.ToQuotedString(v)
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
	hist := string(sequence.ToString(seq))
	_ = r.rl.SaveHistory(hist)
}

func (*sentinel) Equal(data.Value) bool {
	return false
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
		return errors.New(data.ToString(i))
	default:
		panic(debug.ProgrammerError("non-standard error: %s", i))
	}
}

func isRecoverable(err error) bool {
	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		err = e
	}
	msg := err.Error()
	return msg == parse.ErrListNotClosed ||
		msg == parse.ErrVectorNotClosed ||
		msg == parse.ErrObjectNotClosed ||
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
	r.registerBuiltIn("cls", data.MakeProcedure(cls, 0))
	r.registerBuiltIn("doc", doc)
	r.registerBuiltIn("debug", data.MakeProcedure(debugInfo, 0))
	r.registerBuiltIn("help", data.MakeProcedure(help, 0))
	r.registerBuiltIn("quit", data.MakeProcedure(shutdown, 0))
	r.registerBuiltIn("use", r.makeUse())
}

func (r *REPL) registerBuiltIn(n data.Local, v data.Value) {
	ns := r.getBuiltInsNamespace()
	_ = ns.Declare(n).Bind(v)
}

func (r *REPL) makeUse() data.Value {
	return special.Call(func(e encoder.Encoder, args ...data.Value) error {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			return err
		}
		n := args[0].(data.Local)
		old := r.ns
		r.ns = r.ns.Environment().GetQualified(n)
		if old != r.ns {
			fmt.Println()
		}
		return generate.Literal(e, nothing)
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

func formatForREPL(s string) (string, error) {
	md, err := markdown.FormatMarkdown(s)
	if err != nil {
		return "", err
	}
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
	return strings.Join(out, "\n"), nil
}

func help(...data.Value) data.Value {
	md, err := docstring.Get("help")
	if err != nil {
		panic(debug.ProgrammerError(err.Error()))
	}
	out, err := formatForREPL(md)
	if err != nil {
		panic(debug.ProgrammerError(err.Error()))
	}
	fmt.Println(out)
	return nothing
}

var doc = special.Call(func(e encoder.Encoder, args ...data.Value) error {
	if err := data.CheckRangedArity(0, 1, len(args)); err != nil {
		return err
	}
	if len(args) == 0 {
		docSymbolList()
	} else {
		docSymbol(args[0].(data.Local))
	}
	return generate.Literal(e, nothing)
})

func docSymbol(sym data.Symbol) {
	name := string(sym.Name())
	if name == "doc" {
		docSymbolList()
		return
	}
	docStr := docstring.MustGet(name)
	out, err := formatForREPL(docStr)
	if err != nil {
		panic(debug.ProgrammerError(err.Error()))
	}
	fmt.Println(out)
}

func docSymbolList() {
	names := docstring.Names()
	names = escapeNames(names)
	joined := strings.Join(names, ", ")
	out, err := formatForREPL(fmt.Sprintf(docTemplate, joined))
	if err != nil {
		panic(debug.ProgrammerError(err.Error()))
	}
	fmt.Println(out)
}

func makeUserNamespace() env.Namespace {
	return bootstrap.TopLevelEnvironment().GetQualified(UserDomain)
}

func getHistoryFile() string {
	if usr, err := user.Current(); err == nil {
		return path.Join(usr.HomeDir, ".ale-history")
	}
	return ""
}

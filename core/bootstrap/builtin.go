package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/macro"
)

// Error messages
const (
	ErrBuiltInNotFound = "built-in not found: %s"
	ErrSpecialNotFound = "special form not found: %s"
	ErrMacroNotFound   = "macro not found: %s"
)

const (
	defBuiltInName = "def-builtin"
	defSpecialName = "def-special"
	defMacroName   = "def-macro"
)

func (b *bootstrap) builtIns() {
	b.initialFunctions()
	b.specialForms()
	b.availableFunctions()
}

func (b *bootstrap) initialFunctions() {
	ns := b.environment.GetRoot()

	ns.Private(defBuiltInName).Bind(
		makeDefiner(b.funcMap, ErrBuiltInNotFound),
	)
	ns.Private(defSpecialName).Bind(
		makeDefiner(b.specialMap, ErrSpecialNotFound),
	)
	ns.Private(defMacroName).Bind(
		makeDefiner(b.macroMap, ErrMacroNotFound),
	)
}

func (b *bootstrap) specialForms() {
	b.specials(map[data.Local]special.Call{
		"asm*":          builtin.Asm,
		"begin":         builtin.Begin,
		"eval":          builtin.Eval,
		"lambda":        builtin.Lambda,
		"let":           builtin.Let,
		"let-rec":       builtin.LetMutual,
		"macroexpand-1": builtin.MacroExpand1,
		"macroexpand":   builtin.MacroExpand,
	})
}

func (b *bootstrap) availableFunctions() {
	b.functions(map[data.Local]data.Lambda{
		"append":       builtin.Append,
		"assoc":        builtin.Assoc,
		"chan":         builtin.Chan,
		"current-time": builtin.CurrentTime,
		"defer*":       builtin.Defer,
		"dissoc":       builtin.Dissoc,
		"promise*":     builtin.Promise,
		"gensym":       builtin.GenSym,
		"get":          builtin.Get,
		"go*":          builtin.Go,
		"is-a*":        builtin.IsA,
		"lazy-seq*":    builtin.LazySequence,
		"length":       builtin.Length,
		"list":         builtin.List,
		"macro*":       builtin.Macro,
		"nth":          builtin.Nth,
		"object":       builtin.Object,
		"read":         builtin.Read,
		"recover":      builtin.Recover,
		"reverse":      builtin.Reverse,
		"str!":         builtin.ReaderStr,
		"str":          builtin.Str,
		"sym":          builtin.Sym,
		"type-of*":     builtin.TypeOf,
		"vector":       builtin.Vector,
	})

	b.macros(map[data.Local]macro.Call{
		"syntax-quote": macro.SyntaxQuote,
	})
}

func (b *bootstrap) functions(f map[data.Local]data.Lambda) {
	for k, v := range f {
		b.function(k, v)
	}
}

func (b *bootstrap) function(name data.Local, call data.Lambda) {
	b.funcMap[name] = call
}

func (b *bootstrap) macros(m map[data.Local]macro.Call) {
	for k, v := range m {
		b.macro(k, v)
	}
}
func (b *bootstrap) macro(name data.Local, call macro.Call) {
	b.macroMap[name] = call
}

func (b *bootstrap) specials(s map[data.Local]special.Call) {
	for k, v := range s {
		b.special(k, v)
	}
}

func (b *bootstrap) special(name data.Local, call special.Call) {
	b.specialMap[name] = call
}

func makeDefiner[T data.Value](m map[data.Local]T, err string) special.Call {
	return func(e encoder.Encoder, args ...data.Value) {
		data.AssertFixed(1, len(args))
		n := args[0].(data.Local)
		if sf, ok := m[n]; ok {
			e.Globals().Declare(n).Bind(sf)
			generate.Symbol(e, n)
			return
		}
		panic(fmt.Errorf(err, n))
	}
}

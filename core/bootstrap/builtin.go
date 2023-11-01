package bootstrap

import (
	"fmt"

	builtin2 "github.com/kode4food/ale/core/builtin"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
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
		makeDefiner(ns, b.funcMap, ErrBuiltInNotFound),
	)
	ns.Private(defSpecialName).Bind(
		makeDefiner(ns, b.specialMap, ErrSpecialNotFound),
	)
	ns.Private(defMacroName).Bind(
		makeDefiner(ns, b.macroMap, ErrMacroNotFound),
	)
}

func (b *bootstrap) specialForms() {
	b.specials(map[data.Local]special.Call{
		"asm*":          builtin2.Asm,
		"begin":         builtin2.Begin,
		"eval":          builtin2.Eval,
		"lambda":        builtin2.Lambda,
		"let":           builtin2.Let,
		"let-rec":       builtin2.LetMutual,
		"macroexpand-1": builtin2.MacroExpand1,
		"macroexpand":   builtin2.MacroExpand,
	})
}

func (b *bootstrap) availableFunctions() {
	b.functions(map[data.Local]data.Lambda{
		"append":       builtin2.Append,
		"assoc":        builtin2.Assoc,
		"chan":         builtin2.Chan,
		"current-time": builtin2.CurrentTime,
		"defer*":       builtin2.Defer,
		"dissoc":       builtin2.Dissoc,
		"promise*":     builtin2.Promise,
		"gensym":       builtin2.GenSym,
		"get":          builtin2.Get,
		"go*":          builtin2.Go,
		"is-a*":        builtin2.IsA,
		"lazy-seq*":    builtin2.LazySequence,
		"length":       builtin2.Length,
		"list":         builtin2.List,
		"macro*":       builtin2.Macro,
		"nth":          builtin2.Nth,
		"object":       builtin2.Object,
		"read":         builtin2.Read,
		"recover":      builtin2.Recover,
		"reverse":      builtin2.Reverse,
		"str!":         builtin2.ReaderStr,
		"str":          builtin2.Str,
		"sym":          builtin2.Sym,
		"type-of*":     builtin2.TypeOf,
		"vector":       builtin2.Vector,
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

func makeDefiner[T data.Value](
	ns env.Namespace, m map[data.Local]T, err string,
) special.Call {
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

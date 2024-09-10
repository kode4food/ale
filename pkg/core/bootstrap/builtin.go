package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/core/asm"
	coreSpecial "github.com/kode4food/ale/pkg/core/special"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/macro"
)

const (
	// ErrBuiltInNotFound is raised when a built-in procedure can't be resolved
	ErrBuiltInNotFound = "built-in not found: %s"

	// ErrSpecialNotFound is raised when a built-in special can't be resolved
	ErrSpecialNotFound = "special form not found: %s"

	// ErrMacroNotFound is raised when a built-in macro can't be resolved
	ErrMacroNotFound = "macro not found: %s"
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
		makeDefiner(b.procMap, ErrBuiltInNotFound),
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
		"asm*":          asm.Asm,
		"eval":          coreSpecial.Eval,
		"lambda":        coreSpecial.Lambda,
		"let":           coreSpecial.Let,
		"let-rec":       coreSpecial.LetMutual,
		"macroexpand-1": coreSpecial.MacroExpand1,
		"macroexpand":   coreSpecial.MacroExpand,
	})
}

func (b *bootstrap) availableFunctions() {
	b.functions(map[data.Local]data.Procedure{
		"append":       core.Append,
		"assoc":        core.Assoc,
		"chan":         core.Chan,
		"current-time": core.CurrentTime,
		"defer*":       core.Defer,
		"dissoc":       core.Dissoc,
		"promise*":     core.Promise,
		"gensym":       core.GenSym,
		"get":          core.Get,
		"go*":          core.Go,
		"is-a*":        core.IsA,
		"lazy-seq*":    core.LazySequence,
		"length":       core.Length,
		"list":         core.List,
		"macro*":       core.Macro,
		"nth":          core.Nth,
		"object":       core.Object,
		"read":         core.Read,
		"recover":      core.Recover,
		"reverse":      core.Reverse,
		"str!":         core.ReaderStr,
		"str":          core.Str,
		"sym":          core.Sym,
		"type-of*":     core.TypeOf,
		"vector":       core.Vector,
	})

	b.macros(map[data.Local]macro.Call{
		"syntax-quote": macro.SyntaxQuote,
	})
}

func (b *bootstrap) functions(f map[data.Local]data.Procedure) {
	for k, v := range f {
		b.function(k, v)
	}
}

func (b *bootstrap) function(name data.Local, call data.Procedure) {
	b.procMap[name] = call
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

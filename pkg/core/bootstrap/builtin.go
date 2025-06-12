package bootstrap

import (
	"github.com/kode4food/ale/pkg/core/builtin"
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

func (b *bootstrap) populateBuiltins() {
	b.functions(map[data.Local]data.Procedure{
		"append":       builtin.Append,
		"assoc":        builtin.Assoc,
		"chan":         builtin.Chan,
		"current-time": builtin.CurrentTime,
		"defer*":       builtin.Defer,
		"dissoc":       builtin.Dissoc,
		"delay*":       builtin.Delay,
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
		"#include":     builtin.Include,
		"syntax-quote": builtin.SyntaxQuote,
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

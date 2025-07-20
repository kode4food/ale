package bootstrap

import (
	"github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/macro"
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
		env.Bytes:       builtin.Bytes,
		env.Chan:        builtin.Chan,
		env.CurrentTime: builtin.CurrentTime,
		env.Defer:       builtin.Defer,
		env.GenSym:      builtin.GenSym,
		env.Go:          builtin.Go,
		env.IsA:         builtin.IsA,
		env.List:        builtin.List,
		env.Macro:       builtin.Macro,
		env.Object:      builtin.Object,
		env.Read:        builtin.Read,
		env.Recover:     builtin.Recover,
		env.ReaderStr:   builtin.ReaderStr,
		env.Str:         builtin.Str,
		env.Sym:         builtin.Sym,
		env.TypeOf:      builtin.TypeOf,
		env.Vector:      builtin.Vector,
	})

	b.macros(map[data.Local]macro.Call{
		env.SyntaxQuote: builtin.SyntaxQuote,
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

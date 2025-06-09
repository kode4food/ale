package bootstrap

import (
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/pkg/core/special"
	"github.com/kode4food/ale/pkg/data"
)

func (b *bootstrap) populateSpecialForms() {
	b.specials(map[data.Local]compiler.Call{
		"asm*":          special.Asm,
		"eval":          special.Eval,
		"import":        special.Import,
		"lambda":        special.Lambda,
		"let":           special.Let,
		"let-rec":       special.LetMutual,
		"macroexpand-1": special.MacroExpand1,
		"macroexpand":   special.MacroExpand,
		"special*":      special.Special,

		"declared*": special.Declared,
		"in-ns*":    special.InNamespace,
	})
}

func (b *bootstrap) specials(s map[data.Local]compiler.Call) {
	for k, v := range s {
		b.special(k, v)
	}
}

func (b *bootstrap) special(name data.Local, call compiler.Call) {
	b.specialMap[name] = call
}

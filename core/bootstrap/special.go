package bootstrap

import (
	"github.com/kode4food/ale/core/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/lang/env"
)

func (b *bootstrap) populateSpecialForms() {
	b.specials(map[data.Local]compiler.Call{
		env.Asm:          special.Asm,
		env.Eval:         special.Eval,
		env.Import:       special.Import,
		env.Lambda:       special.Lambda,
		env.Let:          special.Let,
		env.LetMutual:    special.LetMutual,
		env.MacroExpand1: special.MacroExpand1,
		env.MacroExpand:  special.MacroExpand,
		env.Special:      special.Special,

		env.Declared:      special.Declared,
		env.MakeNamespace: special.MakeNamespace,
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

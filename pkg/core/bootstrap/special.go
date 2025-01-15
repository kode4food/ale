package bootstrap

import (
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/pkg/core/asm"
	coreSpecial "github.com/kode4food/ale/pkg/core/special"
	"github.com/kode4food/ale/pkg/data"
)

func (b *bootstrap) populateSpecialForms() {
	b.specials(map[data.Local]special.Call{
		"asm*":          asm.Asm,
		"eval":          coreSpecial.Eval,
		"lambda":        coreSpecial.Lambda,
		"let":           coreSpecial.Let,
		"let-rec":       coreSpecial.LetMutual,
		"macroexpand-1": coreSpecial.MacroExpand1,
		"macroexpand":   coreSpecial.MacroExpand,
		"special*":      asm.Special,
	})
}

func (b *bootstrap) specials(s map[data.Local]special.Call) {
	for k, v := range s {
		b.special(k, v)
	}
}

func (b *bootstrap) special(name data.Local, call special.Call) {
	b.specialMap[name] = call
}

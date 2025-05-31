package bootstrap

import (
	"github.com/kode4food/ale/internal/compiler/special"
	core "github.com/kode4food/ale/pkg/core/special"
	"github.com/kode4food/ale/pkg/data"
)

func (b *bootstrap) populateSpecialForms() {
	b.specials(map[data.Local]special.Call{
		"asm*":          core.Asm,
		"eval":          core.Eval,
		"lambda":        core.Lambda,
		"let":           core.Let,
		"let-rec":       core.LetMutual,
		"macroexpand-1": core.MacroExpand1,
		"macroexpand":   core.MacroExpand,
		"special*":      core.Special,
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

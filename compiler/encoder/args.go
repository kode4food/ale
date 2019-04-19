package encoder

import (
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/runtime/isa"
)

type (
	argsInfo struct {
		names data.Names
		rest  bool
	}

	argsStack []*argsInfo
)

func (e *encoder) PushArgs(names data.Names, rest bool) {
	e.args = append(e.args, &argsInfo{
		names: names,
		rest:  rest,
	})
}

func (e *encoder) PopArgs() {
	args := e.args
	al := len(args)
	e.args = args[0 : al-1]
}

func (e *encoder) ResolveArg(l data.LocalSymbol) (isa.Index, bool, bool) {
	lookup := l.Name()
	args := e.args
	for i := len(args) - 1; i >= 0; i-- {
		a := args[i]
		if idx, rest, ok := a.resolveArg(lookup); ok {
			return idx, rest, ok
		}
	}
	return 0, false, false
}

func (a *argsInfo) resolveArg(lookup data.Name) (isa.Index, bool, bool) {
	for idx, n := range a.names {
		if n == lookup {
			nl := len(a.names)
			isRest := a.rest && idx == nl-1
			return isa.Index(idx), isRest, true
		}
	}
	return 0, false, false
}

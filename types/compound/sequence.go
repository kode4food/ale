package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

// List declares a new PairType that will only allow a List with
// elements of the provided elem Type
func List(elem types.Type) types.Type {
	return makeSequence(basic.List, elem)
}

func makeSequence(base, elem types.Type) types.Type {
	name := fmt.Sprintf("%s(%s)", base.Name(), elem.Name())
	ext := extended.New(base)
	first := &namedPair{
		name: name,
		pair: &pair{
			Extended: ext,
			car:      elem,
		},
	}
	rest := &namedUnion{
		name: name,
		union: &union{
			Extended: ext,
			options:  typeList{first, basic.Null},
		},
	}
	first.cdr = rest
	return rest
}

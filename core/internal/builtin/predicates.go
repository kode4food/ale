package builtin

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/async"
	"github.com/kode4food/ale/macro"
)

type predicate func(data.Value) bool

// Error messages
const (
	ErrUnknownPredicate = "unknown predicate: %s"
)

var predicates = map[data.Keyword]predicate{
	"atom":      isAtom,
	"appender":  makeGoTypePredicate[data.Appender](),
	"boolean":   makeGoTypePredicate[data.Bool](),
	"cons":      makeGoTypePredicate[*data.Cons](),
	"counted":   makeGoTypePredicate[data.Counted](),
	"func":      makeGoTypePredicate[data.Function](),
	"indexed":   makeGoTypePredicate[data.Indexed](),
	"keyword":   makeGoTypePredicate[data.Keyword](),
	"list":      makeGoTypePredicate[data.List](),
	"local":     makeGoTypePredicate[data.Local](),
	"macro":     makeGoTypePredicate[macro.Call](),
	"mapped":    makeGoTypePredicate[data.Mapper](),
	"nan":       isNaN,
	"number":    makeGoTypePredicate[data.Number](),
	"object":    makeGoTypePredicate[data.Object](),
	"pair":      makeGoTypePredicate[data.Pair](),
	"promise":   makeGoTypePredicate[async.Promise](),
	"qualified": makeGoTypePredicate[data.Qualified](),
	"resolved":  isResolved,
	"reverser":  makeGoTypePredicate[data.Reverser](),
	"seq":       makeGoTypePredicate[data.Sequence](),
	"special":   makeGoTypePredicate[encoder.Call](),
	"string":    makeGoTypePredicate[data.String](),
	"symbol":    makeGoTypePredicate[data.Symbol](),
	"vector":    makeGoTypePredicate[data.Vector](),
}

// IsA will allow for a little more flexibility in type checking
var IsA = data.Applicative(func(args ...data.Value) data.Value {
	kwd := args[0].(data.Keyword)
	p, ok := predicates[kwd]
	if !ok {
		panic(fmt.Sprintf(ErrUnknownPredicate, kwd))
	}
	if len(args) == 2 {
		return data.Bool(p(args[1]))
	}
	return makePredicate(p)
}, 1, 2)

func makeGoTypePredicate[T any]() predicate {
	return func(v data.Value) bool {
		_, ok := v.(T)
		return ok
	}
}

func makePredicate(p predicate) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		return data.Bool(p(args[0]))
	}, 1)
}

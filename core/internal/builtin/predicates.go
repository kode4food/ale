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

const (
	AtomKey      = data.Keyword("atom")
	AppenderKey  = data.Keyword("appender")
	BooleanKey   = data.Keyword("boolean")
	ConsKey      = data.Keyword("cons")
	CountedKey   = data.Keyword("counted")
	FuncKey      = data.Keyword("func")
	IndexedKey   = data.Keyword("indexed")
	KeywordKey   = data.Keyword("keyword")
	ListKey      = data.Keyword("list")
	LocalKey     = data.Keyword("local")
	MacroKey     = data.Keyword("macro")
	MappedKey    = data.Keyword("mapped")
	NaNKey       = data.Keyword("nan")
	NumberKey    = data.Keyword("number")
	ObjectKey    = data.Keyword("object")
	PairKey      = data.Keyword("pair")
	PromiseKey   = data.Keyword("promise")
	QualifiedKey = data.Keyword("qualified")
	ResolvedKey  = data.Keyword("resolved")
	ReverserKey  = data.Keyword("reverser")
	SeqKey       = data.Keyword("seq")
	SpecialKey   = data.Keyword("special")
	StringKey    = data.Keyword("string")
	SymbolKey    = data.Keyword("symbol")
	VectorKey    = data.Keyword("vector")
)

var predicates = map[data.Keyword]predicate{
	AtomKey:      isAtom,
	AppenderKey:  makeGoTypePredicate[data.Appender](),
	BooleanKey:   makeGoTypePredicate[data.Bool](),
	ConsKey:      makeGoTypePredicate[*data.Cons](),
	CountedKey:   makeGoTypePredicate[data.Counted](),
	FuncKey:      makeGoTypePredicate[data.Function](),
	IndexedKey:   makeGoTypePredicate[data.Indexed](),
	KeywordKey:   makeGoTypePredicate[data.Keyword](),
	ListKey:      makeGoTypePredicate[data.List](),
	LocalKey:     makeGoTypePredicate[data.Local](),
	MacroKey:     makeGoTypePredicate[macro.Call](),
	MappedKey:    makeGoTypePredicate[data.Mapper](),
	NaNKey:       isNaN,
	NumberKey:    makeGoTypePredicate[data.Number](),
	ObjectKey:    makeGoTypePredicate[data.Object](),
	PairKey:      makeGoTypePredicate[data.Pair](),
	PromiseKey:   makeGoTypePredicate[async.Promise](),
	QualifiedKey: makeGoTypePredicate[data.Qualified](),
	ResolvedKey:  isResolved,
	ReverserKey:  makeGoTypePredicate[data.Reverser](),
	SeqKey:       makeGoTypePredicate[data.Sequence](),
	SpecialKey:   makeGoTypePredicate[encoder.Call](),
	StringKey:    makeGoTypePredicate[data.String](),
	SymbolKey:    makeGoTypePredicate[data.Symbol](),
	VectorKey:    makeGoTypePredicate[data.Vector](),
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

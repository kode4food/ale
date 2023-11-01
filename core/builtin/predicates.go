package builtin

import (
	"fmt"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/async"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/macro"
)

type predicate func(data.Value) bool

// Error messages
const (
	ErrUnknownPredicate = "unknown predicate: %s"
)

const (
	AnyKey       = data.Keyword("any")
	AtomKey      = data.Keyword("atom")
	AppenderKey  = data.Keyword("appender")
	BooleanKey   = data.Keyword("boolean")
	ConsKey      = data.Keyword("cons")
	CountedKey   = data.Keyword("counted")
	FunctionKey  = data.Keyword("function")
	IndexedKey   = data.Keyword("indexed")
	KeywordKey   = data.Keyword("keyword")
	ListKey      = data.Keyword("list")
	LocalKey     = data.Keyword("local")
	MacroKey     = data.Keyword("macro")
	MappedKey    = data.Keyword("mapped")
	NaNKey       = data.Keyword("nan")
	NullKey      = data.Keyword("null")
	NumberKey    = data.Keyword("number")
	ObjectKey    = data.Keyword("object")
	PairKey      = data.Keyword("pair")
	PromiseKey   = data.Keyword("promise")
	QualifiedKey = data.Keyword("qualified")
	ResolvedKey  = data.Keyword("resolved")
	ReverserKey  = data.Keyword("reverser")
	SequenceKey  = data.Keyword("sequence")
	SpecialKey   = data.Keyword("special")
	StringKey    = data.Keyword("string")
	SymbolKey    = data.Keyword("symbol")
	VectorKey    = data.Keyword("vector")
)

var listType = types.MakeUnion(types.BasicList, types.BasicNull)

var predicates = map[data.Keyword]data.Lambda{
	AtomKey:     makePredicate(isAtom),
	NaNKey:      makePredicate(isNaN),
	PairKey:     makePredicate(isPair),
	ResolvedKey: makePredicate(isResolved),

	AnyKey:      data.MakeTypePredicate(types.BasicAny),
	BooleanKey:  data.MakeTypePredicate(types.BasicBoolean),
	ConsKey:     data.MakeTypePredicate(types.BasicCons),
	FunctionKey: data.MakeTypePredicate(types.BasicLambda),
	KeywordKey:  data.MakeTypePredicate(types.BasicKeyword),
	ListKey:     data.MakeTypePredicate(listType),
	MacroKey:    data.MakeTypePredicate(macro.CallType),
	NullKey:     data.MakeTypePredicate(types.BasicNull),
	NumberKey:   data.MakeTypePredicate(types.BasicNumber),
	ObjectKey:   data.MakeTypePredicate(types.BasicObject),
	PromiseKey:  data.MakeTypePredicate(async.PromiseType),
	SpecialKey:  data.MakeTypePredicate(special.CallType),
	StringKey:   data.MakeTypePredicate(types.BasicString),
	SymbolKey:   data.MakeTypePredicate(types.BasicSymbol),
	VectorKey:   data.MakeTypePredicate(types.BasicVector),

	AppenderKey:  makeGoTypePredicate[data.Appender](),
	CountedKey:   makeGoTypePredicate[data.Counted](),
	IndexedKey:   makeGoTypePredicate[data.Indexed](),
	LocalKey:     makeGoTypePredicate[data.Local](),
	MappedKey:    makeGoTypePredicate[data.Mapper](),
	QualifiedKey: makeGoTypePredicate[data.Qualified](),
	ReverserKey:  makeGoTypePredicate[data.Reverser](),
	SequenceKey:  makeGoTypePredicate[data.Sequence](),
}

// TypeOf returns a CallType Predicate for the Types of the given Values. If
// more than one Value is provided, the Union of their Types will be returned
var TypeOf = data.MakeLambda(func(args ...data.Value) data.Value {
	return data.TypePredicateOf(args[0], args[1:]...)
}, 1, data.OrMore)

// IsA returns a Predicate from the set of builtin named Predicates
var IsA = data.MakeLambda(func(args ...data.Value) data.Value {
	kwd := args[0].(data.Keyword)
	if p, ok := predicates[kwd]; ok {
		return p
	}
	panic(fmt.Errorf(ErrUnknownPredicate, kwd))
}, 1)

func makeGoTypePredicate[T any]() data.Lambda {
	return data.MakeLambda(func(args ...data.Value) data.Value {
		_, ok := args[0].(T)
		return data.Bool(ok)
	}, 1)
}

func makePredicate(p predicate) data.Lambda {
	return data.MakeLambda(func(args ...data.Value) data.Value {
		return data.Bool(p(args[0]))
	}, 1)
}

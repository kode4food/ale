package builtin

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/sync"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/macro"
)

type predicate func(ale.Value) bool

// ErrUnknownPredicate is raised when a call to IsA can't resolve a built-in
// predicate for the specified keyword
var ErrUnknownPredicate = errors.New("unknown predicate")

const (
	AnyKey       = data.Keyword("any")
	AtomKey      = data.Keyword("atom")
	AppenderKey  = data.Keyword("appender")
	BooleanKey   = data.Keyword("boolean")
	BytesKey     = data.Keyword("bytes")
	ConsKey      = data.Keyword("cons")
	CountedKey   = data.Keyword("counted")
	ProcedureKey = data.Keyword("procedure")
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

var (
	listType = types.MakeUnion(types.BasicList, types.BasicNull)

	predicates = map[data.Keyword]data.Procedure{
		AtomKey:     makePredicate(isAtom),
		NaNKey:      makePredicate(isNaN),
		PairKey:     makePredicate(isPair),
		ResolvedKey: makePredicate(isResolved),

		AnyKey:       data.MakeTypePredicate(types.BasicAny),
		BooleanKey:   data.MakeTypePredicate(types.BasicBoolean),
		BytesKey:     data.MakeTypePredicate(types.BasicBytes),
		ConsKey:      data.MakeTypePredicate(types.BasicCons),
		ProcedureKey: data.MakeTypePredicate(types.BasicProcedure),
		KeywordKey:   data.MakeTypePredicate(types.BasicKeyword),
		ListKey:      data.MakeTypePredicate(listType),
		MacroKey:     data.MakeTypePredicate(macro.CallType),
		NullKey:      data.MakeTypePredicate(types.BasicNull),
		NumberKey:    data.MakeTypePredicate(types.BasicNumber),
		ObjectKey:    data.MakeTypePredicate(types.BasicObject),
		PromiseKey:   data.MakeTypePredicate(sync.PromiseType),
		SpecialKey:   data.MakeTypePredicate(compiler.CallType),
		StringKey:    data.MakeTypePredicate(types.BasicString),
		SymbolKey:    data.MakeTypePredicate(types.BasicSymbol),
		VectorKey:    data.MakeTypePredicate(types.BasicVector),

		AppenderKey:  makeGoTypePredicate[data.Appender](),
		CountedKey:   makeGoTypePredicate[data.Counted](),
		IndexedKey:   makeGoTypePredicate[data.Indexed](),
		LocalKey:     makeGoTypePredicate[data.Local](),
		MappedKey:    makeGoTypePredicate[data.Mapper](),
		QualifiedKey: makeGoTypePredicate[data.Qualified](),
		ReverserKey:  makeGoTypePredicate[data.Reverser](),
		SequenceKey:  makeGoTypePredicate[data.Sequence](),
	}
)

// TypeOf returns a Type Predicate for the Types of the given Values. If more
// than one Value is provided, the Union of their Types will be returned
var TypeOf = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	return data.TypePredicateOf(args[0], args[1:]...)
}, 1, data.OrMore)

// IsA returns a Predicate from the set of builtin named Predicates
var IsA = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	kwd := args[0].(data.Keyword)
	if p, ok := predicates[kwd]; ok {
		return p
	}
	panic(fmt.Errorf("%w: %s", ErrUnknownPredicate, kwd))
}, 1)

func makeGoTypePredicate[T any]() data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		_, ok := args[0].(T)
		return data.Bool(ok)
	}, 1)
}

func makePredicate(p predicate) data.Procedure {
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
		return data.Bool(p(args[0]))
	}, 1)
}

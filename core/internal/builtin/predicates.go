package builtin

import (
	"github.com/kode4food/ale/compiler"
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/async"
	"github.com/kode4food/ale/macro"
)

// Type-checking predicates
var (
	// IsBoolean returns whether the provided value is a boolean
	IsBoolean = makePredicate[data.Bool]()

	// IsKeyword returns whether the provided value is a keyword
	IsKeyword = makePredicate[data.Keyword]()

	// IsApply tests whether a value is callable
	IsApply = makePredicate[data.Function]()

	// IsSpecial tests whether not a function is a special form
	IsSpecial = makePredicate[encoder.Call]()

	// IsVector returns whether the provided value is a vector
	IsVector = makePredicate[data.Vector]()

	// IsAppender returns whether the provided value is an appender
	IsAppender = makePredicate[data.Appender]()

	// IsPromise returns whether the specified value is a promise
	IsPromise = makePredicate[async.Promise]()

	// IsString returns whether the provided value is a string
	IsString = makePredicate[data.String]()

	// IsMacro returns whether the argument is a macro
	IsMacro = makePredicate[macro.Call]()

	// IsSymbol returns whether the provided value is a symbol
	IsSymbol = makePredicate[data.Symbol]()

	// IsLocal returns whether the provided value is an unqualified symbol
	IsLocal = makePredicate[data.Local]()

	// IsQualified returns whether the provided value is a qualified symbol
	IsQualified = makePredicate[data.Qualified]()

	// IsSeq returns whether the provided value is a sequence
	IsSeq = makePredicate[data.Sequence]()

	// IsCounted returns whether the provided value is a counted sequence
	IsCounted = makePredicate[data.Counted]()

	// IsIndexed returns whether the provided value is an indexed sequence
	IsIndexed = makePredicate[data.Indexed]()

	// IsReverser returns whether the value is a reversible sequence
	IsReverser = makePredicate[data.Reverser]()

	// IsNumber returns true if the provided value is a number
	IsNumber = makePredicate[data.Number]()

	// IsList returns whether the provided value is a list
	IsList = makePredicate[data.List]()

	// IsPair returns whether the provided value is a Pair
	IsPair = makePredicate[data.Pair]()

	// IsCons returns whether the provided value is a Cons cell
	IsCons = makePredicate[*data.Cons]()

	// IsObject returns whether a value is an object
	IsObject = makePredicate[data.Object]()

	// IsMapped returns whether a value is a mapped sequence
	IsMapped = makePredicate[data.Mapper]()
)

// IsAtom returns whether the provided value is atomic
var IsAtom = data.Applicative(func(args ...data.Value) data.Value {
	return data.Bool(!compiler.IsEvaluable(args[0]))
}, 1)

func makePredicate[T any]() data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		_, ok := args[0].(T)
		return data.Bool(ok)
	}, 1)
}

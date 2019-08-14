package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/macro"
)

// Error messages
const (
	BuiltInNotFound = "built-in not found: %s"
	SpecialNotFound = "special form not found: %s"
	MacroNotFound   = "macro not found: %s"
)

const (
	defBuiltInName = "def-builtin"
	defSpecialName = "def-special"
	defMacroName   = "def-macro"

	orMore = -1
)

func (b *bootstrap) builtIns() {
	b.specialForms()
	b.initialFunctions()
	b.availableFunctions()
}

func (b *bootstrap) specialForms() {
	b.special("begin", special.Begin)
	b.special("declare*", special.Declare)
	b.special("define*", special.Define)
	b.special("eval", special.Eval)
	b.special("if", special.If)
	b.special("lambda", special.Lambda)
	b.special("let", special.Let)
	b.special("let-rec", special.LetMutual)
	b.special("macroexpand-1", special.MacroExpand1)
	b.special("macroexpand", special.MacroExpand)
	b.special("quote", special.Quote)
}

func (b *bootstrap) initialFunctions() {
	manager := b.manager

	singleArgChecker := arity.MakeFixedChecker(1)

	defBuiltIn := data.MakeNormal(func(args ...data.Value) data.Value {
		ns := manager.GetRoot()
		n := args[0].(data.LocalSymbol).Name()
		if nf, ok := b.funcMap[n]; ok {
			ns.Declare(n).Bind(nf)
			return args[0]
		}
		panic(fmt.Errorf(BuiltInNotFound, n))
	}, singleArgChecker)

	defSpecial := data.MakeNormal(func(args ...data.Value) data.Value {
		ns := manager.GetRoot()
		n := args[0].(data.LocalSymbol).Name()
		if sf, ok := b.specialMap[n]; ok {
			ns.Declare(n).Bind(sf)
			return args[0]
		}
		panic(fmt.Errorf(SpecialNotFound, n))
	}, singleArgChecker)

	defMacro := data.MakeNormal(func(args ...data.Value) data.Value {
		ns := manager.GetRoot()
		n := args[0].(data.LocalSymbol).Name()
		if sf, ok := b.macroMap[n]; ok {
			ns.Declare(n).Bind(sf)
			return args[0]
		}
		panic(fmt.Errorf(MacroNotFound, n))
	}, singleArgChecker)

	ns := b.manager.GetRoot()
	ns.Declare(defBuiltInName).Bind(defBuiltIn)
	ns.Declare(defSpecialName).Bind(defSpecial)
	ns.Declare(defMacroName).Bind(defMacro)
}

func (b *bootstrap) availableFunctions() {
	b.applicative("-", builtin.Sub, 1, orMore)
	b.applicative("!=", builtin.Neq, 1, orMore)
	b.applicative("*", builtin.Mul)
	b.applicative("/", builtin.Div, 1, orMore)
	b.applicative("+", builtin.Add)
	b.applicative("<", builtin.Lt, 1, orMore)
	b.applicative("<=", builtin.Lte, 1, orMore)
	b.applicative("=", builtin.Eq, 1, orMore)
	b.applicative(">", builtin.Gt, 1, orMore)
	b.applicative(">=", builtin.Gte, 1, orMore)

	b.applicative("append", builtin.Append, 2)
	b.applicative("apply", builtin.Apply, 2, orMore)
	b.applicative("car", builtin.Car, 1)
	b.applicative("cdr", builtin.Cdr, 1)
	b.applicative("chan", builtin.Chan, 0, 1)
	b.applicative("cons", builtin.Cons, 2)
	b.applicative("current-time", builtin.CurrentTime, 0)
	b.applicative("defer", builtin.Defer, 2)
	b.applicative("promise", builtin.Promise, 1)
	b.applicative("deque", builtin.Deque)
	b.applicative("eq", builtin.IsIdentical, 1, orMore)
	b.applicative("first", builtin.First, 1)
	b.applicative("gensym", builtin.GenSym, 0, 1)
	b.applicative("get", builtin.Get, 2)
	b.applicative("go*", builtin.Go, 1)
	b.applicative("lazy-seq*", builtin.LazySequence, 1)
	b.applicative("length", builtin.Length, 1)
	b.applicative("list", builtin.List)
	b.applicative("macro", builtin.Macro, 1)
	b.applicative("mod", builtin.Mod, 1, orMore)
	b.applicative("nth", builtin.Nth, 2, 3)
	b.applicative("object", builtin.Object)
	b.applicative("raise", builtin.Raise, 1)
	b.applicative("read", builtin.Read, 1)
	b.applicative("recover", builtin.Recover, 2)
	b.applicative("rest", builtin.Rest, 1)
	b.applicative("reverse", builtin.Reverse, 1)
	b.applicative("str!", builtin.ReaderStr)
	b.applicative("str", builtin.Str)
	b.applicative("sym", builtin.Sym, 1)
	b.applicative("vector", builtin.Vector)

	b.applicative("is-appender", builtin.IsAppender, 1)
	b.applicative("is-apply", builtin.IsApply, 1)
	b.applicative("is-atom", builtin.IsAtom, 1)
	b.applicative("is-boolean", builtin.IsBoolean, 1)
	b.applicative("is-cons", builtin.IsCons, 1)
	b.applicative("is-counted", builtin.IsCounted, 1)
	b.applicative("is-deque", builtin.IsDeque, 1)
	b.applicative("is-empty", builtin.IsEmpty, 1)
	b.applicative("is-indexed", builtin.IsIndexed, 1)
	b.applicative("is-keyword", builtin.IsKeyword, 1)
	b.applicative("is-list", builtin.IsList, 1)
	b.applicative("is-local", builtin.IsLocal, 1)
	b.applicative("is-macro", builtin.IsMacro, 1)
	b.applicative("is-mapped", builtin.IsMapped, 1)
	b.applicative("is-nan", builtin.IsNaN, 1)
	b.applicative("is-neg-inf", builtin.IsNegInf, 1)
	b.applicative("is-number", builtin.IsNumber, 1)
	b.applicative("is-object", builtin.IsObject, 1)
	b.applicative("is-pair", builtin.IsPair, 1)
	b.applicative("is-pos-inf", builtin.IsPosInf, 1)
	b.applicative("is-promise", builtin.IsPromise, 1)
	b.applicative("is-qualified", builtin.IsQualified, 1)
	b.applicative("is-resolved", builtin.IsResolved, 1)
	b.applicative("is-reversible", builtin.IsReverser, 1)
	b.applicative("is-seq", builtin.IsSeq, 1)
	b.applicative("is-special", builtin.IsSpecial, 1)
	b.applicative("is-string", builtin.IsString, 1)
	b.applicative("is-symbol", builtin.IsSymbol, 1)
	b.applicative("is-vector", builtin.IsVector, 1)

	b.macro("syntax-quote", macro.SyntaxQuote)
}

func (b *bootstrap) applicative(name data.Name, call data.Call, arity ...int) {
	fn := data.Call(call)
	b.builtIn(name, fn, arity...)
}

func (b *bootstrap) macro(name data.Name, call macro.Call) {
	b.macroMap[name] = call
}

func (b *bootstrap) special(name data.Name, call encoder.Call) {
	b.specialMap[name] = call
}

func (b *bootstrap) builtIn(name data.Name, call data.Call, a ...int) {
	fn := data.MakeApplicative(call, arity.MakeChecker(a...))
	b.funcMap[name] = fn
}

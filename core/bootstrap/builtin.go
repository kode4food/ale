package bootstrap

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/macro"
)

// Error messages
const (
	ErrBuiltInNotFound = "built-in not found: %s"
	ErrSpecialNotFound = "special form not found: %s"
	ErrMacroNotFound   = "macro not found: %s"
)

const (
	defBuiltInName = "def-builtin"
	defSpecialName = "def-special"
	defMacroName   = "def-macro"
)

func (b *bootstrap) builtIns() {
	b.initialFunctions()
	b.specialForms()
	b.availableFunctions()
}

func (b *bootstrap) initialFunctions() {
	e := b.environment

	defBuiltIn := data.Normal(func(args ...data.Value) data.Value {
		ns := e.GetRoot()
		n := args[0].(data.LocalSymbol).Name()
		if nf, ok := b.funcMap[n]; ok {
			ns.Declare(n).Bind(nf)
			return args[0]
		}
		panic(fmt.Errorf(ErrBuiltInNotFound, n))
	}, 1)

	defSpecial := data.Normal(func(args ...data.Value) data.Value {
		ns := e.GetRoot()
		n := args[0].(data.LocalSymbol).Name()
		if sf, ok := b.specialMap[n]; ok {
			ns.Declare(n).Bind(sf)
			return args[0]
		}
		panic(fmt.Errorf(ErrSpecialNotFound, n))
	}, 1)

	defMacro := data.Normal(func(args ...data.Value) data.Value {
		ns := e.GetRoot()
		n := args[0].(data.LocalSymbol).Name()
		if sf, ok := b.macroMap[n]; ok {
			ns.Declare(n).Bind(sf)
			return args[0]
		}
		panic(fmt.Errorf(ErrMacroNotFound, n))
	}, 1)

	ns := b.environment.GetRoot()
	ns.Declare(defBuiltInName).Bind(defBuiltIn)
	ns.Declare(defSpecialName).Bind(defSpecial)
	ns.Declare(defMacroName).Bind(defMacro)
}

func (b *bootstrap) specialForms() {
	b.specials(map[data.Name]encoder.Call{
		"begin":         special.Begin,
		"declare*":      special.Declare,
		"define*":       special.Define,
		"eval":          special.Eval,
		"if":            special.If,
		"lambda":        special.Lambda,
		"let":           special.Let,
		"let-rec":       special.LetMutual,
		"macroexpand-1": special.MacroExpand1,
		"macroexpand":   special.MacroExpand,
		"quote":         special.Quote,
		"pattern":       special.Pattern,
	})
}

func (b *bootstrap) availableFunctions() {
	b.functions(map[data.Name]data.Function{
		"-":  builtin.Sub,
		"!=": builtin.Neq,
		"*":  builtin.Mul,
		"/":  builtin.Div,
		"+":  builtin.Add,
		"<":  builtin.Lt,
		"<=": builtin.Lte,
		"=":  builtin.Eq,
		">":  builtin.Gt,
		">=": builtin.Gte,

		"append":       builtin.Append,
		"apply":        builtin.Apply,
		"assoc":        builtin.Assoc,
		"car":          builtin.Car,
		"cdr":          builtin.Cdr,
		"chan":         builtin.Chan,
		"cons":         builtin.Cons,
		"current-time": builtin.CurrentTime,
		"defer":        builtin.Defer,
		"dissoc":       builtin.Dissoc,
		"promise":      builtin.Promise,
		"eq":           builtin.IsIdentical,
		"first":        builtin.First,
		"gensym":       builtin.GenSym,
		"get":          builtin.Get,
		"go*":          builtin.Go,
		"lazy-seq*":    builtin.LazySequence,
		"length":       builtin.Length,
		"list":         builtin.List,
		"macro":        builtin.Macro,
		"mod":          builtin.Mod,
		"nth":          builtin.Nth,
		"object":       builtin.Object,
		"raise":        builtin.Raise,
		"read":         builtin.Read,
		"recover":      builtin.Recover,
		"rest":         builtin.Rest,
		"reverse":      builtin.Reverse,
		"str!":         builtin.ReaderStr,
		"str":          builtin.Str,
		"sym":          builtin.Sym,
		"vector":       builtin.Vector,

		"is-appender":   builtin.IsAppender,
		"is-apply":      builtin.IsApply,
		"is-atom":       builtin.IsAtom,
		"is-boolean":    builtin.IsBoolean,
		"is-cons":       builtin.IsCons,
		"is-counted":    builtin.IsCounted,
		"is-empty":      builtin.IsEmpty,
		"is-indexed":    builtin.IsIndexed,
		"is-keyword":    builtin.IsKeyword,
		"is-list":       builtin.IsList,
		"is-local":      builtin.IsLocal,
		"is-macro":      builtin.IsMacro,
		"is-mapped":     builtin.IsMapped,
		"is-nan":        builtin.IsNaN,
		"is-neg-inf":    builtin.IsNegInf,
		"is-number":     builtin.IsNumber,
		"is-object":     builtin.IsObject,
		"is-pair":       builtin.IsPair,
		"is-pos-inf":    builtin.IsPosInf,
		"is-promise":    builtin.IsPromise,
		"is-qualified":  builtin.IsQualified,
		"is-resolved":   builtin.IsResolved,
		"is-reversible": builtin.IsReverser,
		"is-seq":        builtin.IsSeq,
		"is-special":    builtin.IsSpecial,
		"is-string":     builtin.IsString,
		"is-symbol":     builtin.IsSymbol,
		"is-vector":     builtin.IsVector,
	})

	b.macros(map[data.Name]macro.Call{
		"syntax-quote": macro.SyntaxQuote,
	})
}

func (b *bootstrap) functions(f map[data.Name]data.Function) {
	for k, v := range f {
		b.function(k, v)
	}
}

func (b *bootstrap) function(name data.Name, call data.Function) {
	b.funcMap[name] = call
}

func (b *bootstrap) macros(m map[data.Name]macro.Call) {
	for k, v := range m {
		b.macro(k, v)
	}
}
func (b *bootstrap) macro(name data.Name, call macro.Call) {
	b.macroMap[name] = call
}

func (b *bootstrap) specials(s map[data.Name]encoder.Call) {
	for k, v := range s {
		b.special(k, v)
	}
}

func (b *bootstrap) special(name data.Name, call encoder.Call) {
	b.specialMap[name] = call
}

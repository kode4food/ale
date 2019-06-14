package macro

import (
	"fmt"
	"strings"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
)

type syntaxEnv struct {
	namespace namespace.Type
	genSyms   map[string]data.Symbol
}

// Error messages
const (
	UnsupportedSyntaxQuote = "unsupported type in syntax quote: %s"
)

var (
	quoteSym  = namespace.RootSymbol("quote")
	listSym   = namespace.RootSymbol("list")
	vectorSym = namespace.RootSymbol("vector")
	objectSym = namespace.RootSymbol("object")
	applySym  = namespace.RootSymbol("apply")
	concatSym = namespace.RootSymbol("concat!")

	unquoteSym  = namespace.RootSymbol("unquote")
	splicingSym = namespace.RootSymbol("unquote-splicing")
)

// SyntaxQuote performs syntax quoting on the provided value
func SyntaxQuote(ns namespace.Type, args ...data.Value) data.Value {
	arity.AssertFixed(1, len(args))
	value := args[0]
	sc := &syntaxEnv{
		namespace: ns,
		genSyms:   make(map[string]data.Symbol),
	}
	return sc.quote(value)
}

func (se *syntaxEnv) quote(v data.Value) data.Value {
	return se.quoteValue(v)
}

func (se *syntaxEnv) quoteValue(v data.Value) data.Value {
	switch typed := v.(type) {
	case data.Sequence:
		return se.quoteSequence(typed)
	case data.Symbol:
		return se.quoteSymbol(typed)
	default:
		return v
	}
}

func (se *syntaxEnv) quoteSymbol(s data.Symbol) data.Value {
	if gs, ok := se.generateSymbol(s); ok {
		return data.NewList(quoteSym, gs)
	}
	return data.NewList(quoteSym, se.qualifySymbol(s))
}

func (se *syntaxEnv) generateSymbol(s data.Symbol) (data.Symbol, bool) {
	if _, ok := s.(data.QualifiedSymbol); ok {
		return nil, false
	}

	n := string(s.Name())
	if len(n) <= 1 || !strings.HasSuffix(n, "#") {
		return nil, false
	}

	if r, ok := se.genSyms[n]; ok {
		return r, true
	}

	r := data.NewGeneratedSymbol(data.Name(n[0 : len(n)-1]))
	se.genSyms[n] = r
	return r, true
}

func (se *syntaxEnv) quoteSequence(s data.Sequence) data.Value {
	switch typed := s.(type) {
	case data.String:
		return typed
	case data.List:
		return data.NewList(applySym, listSym, se.quoteElements(typed))
	case data.Vector:
		return data.NewList(applySym, vectorSym, se.quoteElements(typed))
	case data.Object:
		return se.quoteObject(typed)
	case data.NullType:
		return typed
	default:
		panic(fmt.Errorf(UnsupportedSyntaxQuote, s))
	}
}

func (se *syntaxEnv) quoteObject(as data.Object) data.Value {
	res := data.EmptyVector
	for f, r, ok := as.Split(); ok; f, r, ok = r.Split() {
		p := f.(data.Pair)
		res = append(res, p.Car(), p.Cdr())
	}
	return data.NewList(applySym, objectSym, se.quoteElements(res))
}

func (se *syntaxEnv) quoteElements(s data.Sequence) data.Value {
	res := data.EmptyVector
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if v, ok := isUnquoteSplicing(f); ok {
			res = append(res, v)
			continue
		}
		if v, ok := isUnquote(f); ok {
			res = append(res, data.NewList(listSym, v))
			continue
		}
		res = append(res, data.NewList(listSym, se.quoteValue(f)))
	}
	return data.NewList(res...).Prepend(concatSym)
}

func (se *syntaxEnv) qualifySymbol(s data.Symbol) data.Value {
	if q, ok := s.(data.QualifiedSymbol); ok {
		return q
	}
	name := s.Name()
	if e, ok := se.namespace.Resolve(name); ok {
		return data.NewQualifiedSymbol(name, e.Owner().Domain())
	}
	return s
}

func isWrapperCall(s data.Symbol, v data.Value) (data.Value, bool) {
	if l, ok := isBuiltInCall(s, v); ok {
		return l.Rest().First(), true
	}
	return data.Null, false
}

func isBuiltInCall(s data.Symbol, v data.Value) (data.List, bool) {
	if l, ok := v.(data.List); ok && l.Count() > 0 {
		if call, ok := l.First().(data.Symbol); ok {
			return l, call == s
		}
	}
	return nil, false
}

func isUnquote(v data.Value) (data.Value, bool) {
	return isWrapperCall(unquoteSym, v)
}

func isUnquoteSplicing(v data.Value) (data.Value, bool) {
	return isWrapperCall(splicingSym, v)
}

package builtin

import (
	"fmt"
	"strings"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/namespace"
)

type syntaxEnv struct {
	genSyms map[string]api.Symbol
}

const (
	// UnsupportedSyntaxQuote is raised when something can't be quoted
	UnsupportedSyntaxQuote = "unsupported type in syntax quote: %s"

	quoteName           = "quote"
	listName            = "list"
	vectorName          = "vector"
	assocName           = "assoc"
	applyName           = "apply"
	concatName          = "concat"
	unquoteName         = "unquote"
	unquoteSplicingName = "unquote-splicing"
)

var (
	quoteSym  = api.NewQualifiedSymbol(quoteName, namespace.RootDomain)
	listSym   = api.NewQualifiedSymbol(listName, namespace.RootDomain)
	vectorSym = api.NewQualifiedSymbol(vectorName, namespace.RootDomain)
	assocSym  = api.NewQualifiedSymbol(assocName, namespace.RootDomain)
	applySym  = api.NewQualifiedSymbol(applyName, namespace.RootDomain)
	concatSym = api.NewQualifiedSymbol(concatName, namespace.RootDomain)
)

// Quote returns the provided value without evaluating it
func Quote(args ...api.Value) api.Value {
	return args[0]
}

// SyntaxQuote performs syntax quoting on the provided value
func SyntaxQuote(args ...api.Value) api.Value {
	sc := &syntaxEnv{
		genSyms: make(map[string]api.Symbol),
	}
	return sc.quote(args[0])
}

func (se *syntaxEnv) quote(v api.Value) api.Value {
	return se.quoteValue(v)
}

func (se *syntaxEnv) quoteValue(v api.Value) api.Value {
	switch typed := v.(type) {
	case api.Sequence:
		return se.quoteSequence(typed)
	case api.Symbol:
		return se.quoteSymbol(typed)
	default:
		return v
	}
}

func (se *syntaxEnv) quoteSymbol(s api.Symbol) api.Value {
	if gs, ok := se.generateSymbol(s); ok {
		return api.NewList(quoteSym, gs)
	}
	return api.NewList(quoteSym, s)
}

func (se *syntaxEnv) generateSymbol(s api.Symbol) (api.Symbol, bool) {
	if _, ok := s.(api.QualifiedSymbol); ok {
		return nil, false
	}

	n := string(s.Name())
	if len(n) <= 1 || !strings.HasSuffix(n, "#") {
		return nil, false
	}

	if r, ok := se.genSyms[n]; ok {
		return r, true
	}

	r := api.NewGeneratedSymbol(api.Name(n[0 : len(n)-1]))
	se.genSyms[n] = r
	return r, true
}

func (se *syntaxEnv) quoteSequence(s api.Sequence) api.Value {
	switch typed := s.(type) {
	case api.String:
		return typed
	case *api.List:
		return api.NewList(applySym, listSym, se.quoteElements(typed))
	case api.Vector:
		return api.NewList(applySym, vectorSym, se.quoteElements(typed))
	case api.Associative:
		return se.quoteAssociative(typed)
	default:
		panic(fmt.Errorf(UnsupportedSyntaxQuote, s))
	}
}

func (se *syntaxEnv) quoteAssociative(as api.Associative) api.Value {
	res := api.EmptyVector
	for f, r, ok := as.Split(); ok; f, r, ok = r.Split() {
		p := f.(api.Vector)
		k, _ := p.ElementAt(0)
		v, _ := p.ElementAt(1)
		res = append(res, k)
		res = append(res, v)
	}
	return api.NewList(applySym, assocSym, se.quoteElements(res))
}

func (se *syntaxEnv) quoteElements(s api.Sequence) api.Value {
	res := api.EmptyVector
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if v, ok := isUnquoteSplicing(f); ok {
			res = append(res, v)
			continue
		}
		if v, ok := isUnquote(f); ok {
			res = append(res, api.NewList(listSym, v))
			continue
		}
		res = append(res, api.NewList(listSym, se.quoteValue(f)))
	}
	return api.NewList(res...).Prepend(concatSym)
}

func isWrapperCall(n api.Name, v api.Value) (api.Value, bool) {
	if l, ok := isBuiltInCall(n, v); ok {
		return l.Rest().First(), true
	}
	return api.Nil, false
}

func isBuiltInDomain(s api.Symbol) bool {
	if qs, ok := s.(api.QualifiedSymbol); ok {
		return qs.Domain() == namespace.RootDomain
	}
	return false
}

func isBuiltInCall(n api.Name, v api.Value) (*api.List, bool) {
	if l, ok := v.(*api.List); ok && l.Count() > 0 {
		if s, ok := l.First().(api.Symbol); ok {
			return l, isBuiltInDomain(s) && s.Name() == n
		}
	}
	return nil, false
}

func isUnquote(v api.Value) (api.Value, bool) {
	return isWrapperCall(unquoteName, v)
}

func isUnquoteSplicing(v api.Value) (api.Value, bool) {
	return isWrapperCall(unquoteSplicingName, v)
}

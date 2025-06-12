package builtin

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

type syntaxEnv struct {
	ns      env.Namespace
	genSyms map[string]data.Symbol
}

// ErrUnsupportedSyntaxQuote is raised when an attempt to syntax quote an
// unsupported type is made. Generally on basic sequences are supported
const ErrUnsupportedSyntaxQuote = "unsupported type in syntax quote: %s"

var (
	quoteSym  = env.RootSymbol("quote")
	consSym   = env.RootSymbol("cons")
	listSym   = env.RootSymbol("list")
	vectorSym = env.RootSymbol("vector")
	objectSym = env.RootSymbol("object")
	applySym  = env.RootSymbol("apply")
	concatSym = env.RootSymbol("concat!")

	unquoteSym  = env.RootSymbol("unquote")
	splicingSym = env.RootSymbol("unquote-splicing")
)

// SyntaxQuote performs syntax quoting on the provided value
func SyntaxQuote(ns env.Namespace, args ...data.Value) data.Value {
	data.MustCheckFixedArity(1, len(args))
	value := args[0]
	sc := &syntaxEnv{
		ns:      ns,
		genSyms: map[string]data.Symbol{},
	}
	res, err := sc.quote(value)
	if err != nil {
		panic(err)
	}
	return res
}

func (se *syntaxEnv) quote(v data.Value) (data.Value, error) {
	return se.quoteValue(v)
}

func (se *syntaxEnv) quoteValue(v data.Value) (data.Value, error) {
	switch v := v.(type) {
	case data.Sequence:
		return se.quoteSequence(v)
	case data.Pair:
		return se.quotePair(v)
	case data.Symbol:
		return se.quoteSymbol(v), nil
	default:
		return v, nil
	}
}

func (se *syntaxEnv) quoteSymbol(s data.Symbol) data.Value {
	if gs, ok := se.generateSymbol(s); ok {
		return data.NewList(quoteSym, gs)
	}
	return data.NewList(quoteSym, se.qualifySymbol(s))
}

func (se *syntaxEnv) generateSymbol(s data.Symbol) (data.Symbol, bool) {
	if _, ok := s.(data.Qualified); ok {
		return nil, false
	}

	n := string(s.Name())
	if len(n) <= 1 || !strings.HasSuffix(n, "#") {
		return nil, false
	}

	if r, ok := se.genSyms[n]; ok {
		return r, true
	}

	r := data.NewGeneratedSymbol(data.Local(n[0 : len(n)-1]))
	se.genSyms[n] = r
	return r, true
}

func (se *syntaxEnv) quoteSequence(s data.Sequence) (data.Value, error) {
	switch s := s.(type) {
	case data.String:
		return s, nil
	case *data.List:
		if s == data.Null {
			return s, nil
		}
		e, err := se.quoteElements(s)
		if err != nil {
			return nil, err
		}
		return data.NewList(applySym, listSym, e), nil
	case data.Vector:
		e, err := se.quoteElements(s)
		if err != nil {
			return nil, err
		}
		return data.NewList(applySym, vectorSym, e), nil
	case *data.Object:
		return se.quoteObject(s)
	default:
		return nil, fmt.Errorf(ErrUnsupportedSyntaxQuote, s)
	}
}

func (se *syntaxEnv) quotePair(c data.Pair) (data.Value, error) {
	car, err := se.quoteValue(c.Car())
	if err != nil {
		return nil, err
	}
	cdr, err := se.quoteValue(c.Cdr())
	if err != nil {
		return nil, err
	}
	return data.NewList(consSym, car, cdr), nil
}

func (se *syntaxEnv) quoteObject(as *data.Object) (data.Value, error) {
	var res data.Vector
	for f, r, ok := as.Split(); ok; f, r, ok = r.Split() {
		p := f.(data.Pair)
		res = append(res, p.Car(), p.Cdr())
	}
	e, err := se.quoteElements(res)
	if err != nil {
		return nil, err
	}
	return data.NewList(applySym, objectSym, e), nil
}

func (se *syntaxEnv) quoteElements(s data.Sequence) (data.Value, error) {
	var res data.Vector
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if v, ok := isUnquoteSplicing(f); ok {
			res = append(res, v)
			continue
		}
		if v, ok := isUnquote(f); ok {
			res = append(res, data.NewList(listSym, v))
			continue
		}
		v, err := se.quoteValue(f)
		if err != nil {
			return nil, err
		}
		res = append(res, data.NewList(listSym, v))
	}
	return data.NewList(res...).Prepend(concatSym), nil
}

func (se *syntaxEnv) qualifySymbol(s data.Symbol) data.Value {
	if q, ok := s.(data.Qualified); ok {
		return q
	}
	name := s.Name()
	if _, in, err := se.ns.Resolve(name); err == nil {
		return data.NewQualifiedSymbol(name, in.Domain())
	}
	return s
}

func isWrapperCall(s data.Symbol, v data.Value) (data.Value, bool) {
	if l, ok := isBuiltInCall(s, v); ok {
		return l.Cdr().(data.Pair).Car(), true
	}
	return data.Null, false
}

func isBuiltInCall(s data.Symbol, v data.Value) (*data.List, bool) {
	if l, ok := v.(*data.List); ok && l.Count() > 0 {
		if call, ok := l.Car().(data.Symbol); ok {
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

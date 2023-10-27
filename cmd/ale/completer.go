package main

import (
	"strings"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/read"
	"github.com/kode4food/ale/read/lex"
	"github.com/kode4food/comb/slices"
)

var emptyStrings []string

func (r *REPL) Do(line []rune, pos int) ([][]rune, int) {
	pfx := string(line[:pos])
	buf := r.buf.String() + pfx
	res, off := r.autoComplete(buf)
	needSpace := pos == len(line) || line[pos] != ' '
	return slices.Map(res, func(s string) []rune {
		elem := []rune(s[off:])
		last := len(elem) - 1
		if !needSpace && elem[last] == ' ' {
			elem = elem[:last]
		}
		return elem
	}), 0
}

func (r *REPL) autoComplete(buf string) ([]string, int) {
	src := data.String(buf)
	seq := read.Tokens(src)
	if l, ok := data.Last(seq); ok {
		if l := l.(*lex.Token); ok && l.Type() == lex.Identifier {
			pfx := string(l.Value().(data.String))
			return r.prefixedSymbols(pfx), len(pfx)
		}
	}
	return emptyStrings, 0
}

func (r *REPL) prefixedSymbols(pfx string) []string {
	s, err := data.ParseSymbol(data.String(pfx))
	if err != nil {
		return emptyStrings
	}
	switch s := s.(type) {
	case data.Local:
		return r.prefixedLocals(s)
	case data.Qualified:
		return r.prefixedQualified(s)
	}
	// Programmer error
	panic("unexpected symbol parsing result")
}

func (r *REPL) prefixedLocals(s data.Local) []string {
	name := s.String()
	root := r.ns.Environment().GetRoot()
	current := r.ns

	var res []string
	res = r.prefixedDomains(res, s)
	res = addPrefixed(res, name, root.Declared())
	if current != root {
		res = addPrefixed(res, name, current.Declared())
	}
	return res
}

func (r *REPL) prefixedDomains(res []string, s data.Local) []string {
	name := s.String()
	return append(res, slices.Map(
		slices.Filter(r.ns.Environment().Domains(),
			func(d data.Local) bool {
				return strings.HasPrefix(string(d), name)
			},
		),
		func(d data.Local) string {
			return string(d) + "/"
		},
	)...)
}

func (r *REPL) prefixedQualified(s data.Qualified) []string {
	domain := s.Domain()
	name := s.Name().String()
	ns := r.ns.Environment().GetQualified(s.Domain())
	return slices.Map(
		slices.Filter(ns.Declared(),
			func(n data.Local) bool {
				return strings.HasPrefix(string(n), name)
			},
		),
		func(n data.Local) string {
			qs := data.NewQualifiedSymbol(n, domain)
			return qs.String() + " "
		},
	)
}

func addPrefixed(res []string, pfx string, names data.Locals) []string {
	return append(res, slices.Map(
		slices.Filter(names, func(n data.Local) bool {
			return strings.HasPrefix(string(n), pfx)
		}), func(n data.Local) string {
			return string(n) + " "
		},
	)...)
}

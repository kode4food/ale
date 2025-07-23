package internal

import (
	"slices"
	"strings"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/read"
)

var emptyStrings []string

func (r *REPL) Do(line []rune, pos int) ([][]rune, int) {
	pfx := string(line[:pos])
	buf := r.buf.String() + pfx
	res, off := r.autoComplete(buf)
	needSpace := pos == len(line) || line[pos] != ' '
	return basics.Map(res, func(s string) []rune {
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
	seq := read.MustTokenize(src)
	if l, ok := sequence.Last(seq); ok {
		if l, ok := l.(*lex.Token); ok && l.Type() == lex.Identifier {
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
	panic(debug.ProgrammerError("unexpected symbol parsing result"))
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
	slices.Sort(res)
	return res
}

func (r *REPL) prefixedDomains(res []string, s data.Local) []string {
	name := s.String()
	return append(res, basics.Map(
		basics.Filter(r.ns.Environment().Domains(),
			func(d data.Local) bool {
				return strings.HasPrefix(string(d), name)
			},
		),
		func(d data.Local) string {
			return string(d) + lang.DomainSeparator
		},
	)...)
}

func (r *REPL) prefixedQualified(s data.Qualified) []string {
	domain := s.Domain()
	name := s.Local().String()
	ns := env.MustGetQualified(r.ns.Environment(), s.Domain())
	res := basics.Map(
		basics.Filter(ns.Declared(),
			func(n data.Local) bool {
				return strings.HasPrefix(string(n), name)
			},
		),
		func(n data.Local) string {
			qs := data.NewQualifiedSymbol(n, domain)
			return data.ToString(qs) + lang.Space
		},
	)
	slices.Sort(res)
	return res
}

func addPrefixed(res []string, pfx string, names data.Locals) []string {
	return append(res, basics.Map(
		basics.Filter(names, func(n data.Local) bool {
			return strings.HasPrefix(string(n), pfx)
		}), func(n data.Local) string {
			return string(n) + lang.Space
		},
	)...)
}

package main

import (
	"strings"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/read"
)

var emptyStrings []string

func (r *REPL) Do(line []rune, pos int) ([][]rune, int) {
	pfx := string(line[:pos])
	buf := r.buf.String() + pfx
	res, off := r.autoComplete(buf)
	out := make([][]rune, len(res))
	needSpace := pos == len(line) || line[pos] != ' '
	for i, s := range res {
		elem := []rune(s[off:])
		last := len(elem) - 1
		if !needSpace && elem[last] == ' ' {
			elem = elem[:last]
		}
		out[i] = elem
	}
	return out, 0
}

func (r *REPL) autoComplete(buf string) ([]string, int) {
	src := data.String(buf)
	seq := read.Tokens(src)
	if l, ok := data.Last(seq); ok {
		if l := l.(*read.Token); ok && l.Type() == read.Identifier {
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
	return emptyStrings
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
	for _, d := range r.ns.Environment().Domains() {
		domain := d.String()
		if strings.HasPrefix(domain, name) {
			res = append(res, domain+"/")
		}
	}
	return res
}

func (r *REPL) prefixedQualified(s data.Qualified) []string {
	domain := s.Domain()
	name := s.Name().String()
	ns := r.ns.Environment().GetQualified(s.Domain())
	var res []string
	for _, n := range ns.Declared() {
		str := n.String()
		if strings.HasPrefix(str, name) {
			qs := data.NewQualifiedSymbol(data.Local(str), domain)
			res = append(res, qs.String()+" ")
		}
	}
	return res
}

func addPrefixed(res []string, pfx string, names data.Locals) []string {
	for _, n := range names {
		str := n.String()
		if strings.HasPrefix(str, pfx) {
			res = append(res, str+" ")
		}
	}
	return res
}

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
	for i, s := range res {
		out[i] = []rune(s[off:] + " ")
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
	var res []string
	root := r.ns.Environment().GetRoot()
	current := r.ns
	res = addPrefixed(res, pfx, root.Declared())
	if current != root {
		res = addPrefixed(res, pfx, current.Declared())
	}
	return res
}

func addPrefixed(res []string, pfx string, names data.Locals) []string {
	for _, n := range names {
		str := n.String()
		if strings.HasPrefix(str, pfx) {
			res = append(res, str)
		}
	}
	return res
}

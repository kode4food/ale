package main

import (
	"strings"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/read"
)

var emptyStrings []string

func (r *REPL) Do(line []rune, pos int) ([][]rune, int) {
	pfx := string(line[:pos])
	res, _ := r.autoCompleter(pfx)
	out := make([][]rune, len(res))
	for i, s := range res {
		out[i] = []rune(s)
	}
	return out, 0
}

func (r *REPL) autoCompleter(line string) ([]string, int) {
	src := data.String(r.buf.String() + line)
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
	root := r.ns.Environment().GetRoot().Declared()
	names := r.ns.Declared()
	res := make([]string, 0, len(root)+len(names))
	res = addPrefixed(res, pfx, root)
	res = addPrefixed(res, pfx, names)
	return res
}

func addPrefixed(res []string, pfx string, names data.Locals) []string {
	for _, n := range names {
		if strings.HasPrefix(n.String(), pfx) {
			res = append(res, n.String()[len(pfx):])
		}
	}
	return res
}

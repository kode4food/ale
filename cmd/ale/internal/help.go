package internal

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale/cmd/ale/internal/docstring"
	"github.com/kode4food/ale/cmd/ale/internal/markdown"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

var docTemplate = docstring.MustGet("doc")

var doc = compiler.Call(func(e encoder.Encoder, args ...data.Value) error {
	if err := data.CheckRangedArity(0, 1, len(args)); err != nil {
		return err
	}
	if len(args) == 0 {
		docSymbolList()
	} else {
		docSymbol(args[0].(data.Local))
	}
	return generate.Literal(e, nothing)
})

func help(...data.Value) data.Value {
	md, err := docstring.Get("help")
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	out, err := formatForREPL(md)
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	fmt.Println(out)
	return nothing
}

func docSymbol(sym data.Symbol) {
	name := string(sym.Local())
	if name == "doc" {
		docSymbolList()
		return
	}
	docStr := docstring.MustGet(name)
	out, err := formatForREPL(docStr)
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	fmt.Println(out)
}

func docSymbolList() {
	names := docstring.Names()
	names = escapeNames(names)
	joined := strings.Join(names, ", ")
	out, err := formatForREPL(fmt.Sprintf(docTemplate, joined))
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	fmt.Println(out)
}

func formatForREPL(s string) (string, error) {
	md, err := markdown.FormatMarkdown(s)
	if err != nil {
		return "", err
	}
	lines := strings.Split(md, "\n")
	var out []string
	out = append(out, "")
	for _, l := range lines {
		if isEmptyString(l) {
			out = append(out, l)
		} else {
			out = append(out, "  "+l)
		}
	}
	out = append(out, "")
	return strings.Join(out, "\n"), nil
}

func escapeNames(names []string) []string {
	return basics.Map(names, func(n string) string {
		if strings.Contains("`*_", n[:1]) {
			return `\` + n
		}
		return n
	})
}

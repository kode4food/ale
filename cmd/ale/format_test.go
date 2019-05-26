package main_test

import (
	"testing"

	main "gitlab.com/kode4food/ale/cmd/ale"
	"gitlab.com/kode4food/ale/cmd/ale/docstring"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestFormatMarkdown(t *testing.T) {
	as := assert.New(t)

	s := docstring.Get("if")
	r := data.String(main.FormatMarkdown(s))
	as.NotContains("---", r)
	as.Contains("\x1b[35m\x1b[1mperforms simple branching\x1b[0m\n\n", r)
}

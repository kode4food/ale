package main_test

import (
	"testing"

	main "github.com/kode4food/ale/cmd/ale"
	"github.com/kode4food/ale/cmd/ale/docstring"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestFormatMarkdown(t *testing.T) {
	as := assert.New(t)

	s, err := docstring.Get("if")
	as.NotEmpty(s)
	as.Nil(err)

	r := data.String(main.FormatMarkdown(s))
	as.NotContains("---", r)
	as.Contains("\x1b[35m\x1b[1mperforms simple branching\x1b[0m\n\n", r)
}

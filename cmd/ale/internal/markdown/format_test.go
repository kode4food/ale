//go:build !windows

package markdown_test

import (
	"testing"

	"github.com/kode4food/ale/cmd/ale/internal/docstring"
	"github.com/kode4food/ale/cmd/ale/internal/markdown"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestFormatMarkdown(t *testing.T) {
	as := assert.New(t)

	s, err := docstring.Get("if")
	if as.NoError(err) {
		as.NotEmpty(s)
	}

	res, err := markdown.FormatMarkdown(s)
	if as.NoError(err) {
		r := S(res)
		as.NotContains("---", r)
		as.Contains("\x1b[35m\x1b[1mperforms simple branching\x1b[0m\n\n", r)
	}
}

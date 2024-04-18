//go:build !windows

package markdown_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/markdown"
	"github.com/kode4food/ale/pkg/core/docstring"
)

func TestFormatMarkdown(t *testing.T) {
	as := assert.New(t)

	s, err := docstring.Get("if")
	as.NotEmpty(s)
	as.Nil(err)

	r := S(markdown.FormatMarkdown(s))
	as.NotContains("---", r)
	as.Contains("\x1b[35m\x1b[1mperforms simple branching\x1b[0m\n\n", r)
}

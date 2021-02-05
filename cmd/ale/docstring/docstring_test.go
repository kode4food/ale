package docstring_test

import (
	"testing"

	"github.com/kode4food/ale/cmd/ale/docstring"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestDocString(t *testing.T) {
	as := assert.New(t)

	ifStr, err := docstring.Get("if")
	as.Contains("---", data.String(ifStr))
	as.Nil(err)

	s, err := docstring.Get("no-way-this-exists")
	as.Empty(s)
	as.Errorf(err, docstring.ErrDocNotFound, "no-way-this-exists")
}

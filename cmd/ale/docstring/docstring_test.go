package docstring_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/cmd/ale/docstring"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestDocString(t *testing.T) {
	as := assert.New(t)

	as.True(docstring.Exists("if"))
	as.False(docstring.Exists("no-way-this-exists"))

	ifStr := docstring.Get("if")
	as.Contains("---", data.String(ifStr))

	errStr := fmt.Sprintf(docstring.ErrDocNotFound, "no-way-this-exists")
	defer as.ExpectPanic(errStr)
	docstring.Get("no-way-this-exists")
}

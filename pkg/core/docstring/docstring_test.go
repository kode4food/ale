package docstring_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/docstring"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestDocString(t *testing.T) {
	as := assert.New(t)

	ifStr, err := docstring.Get("if")
	as.Contains("---", S(ifStr))
	as.Nil(err)

	s, err := docstring.Get("no-way-this-exists")
	as.Empty(s)
	as.EqualError(err, fmt.Sprintf(docstring.ErrSymbolNotDocumented, "no-way-this-exists"))
}

func TestDocumentedBuiltinsExist(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	ns.Declare("doc").Bind(data.Null) // special case for REPL

	for _, name := range docstring.Names() {
		d, _ := docstring.Get(name)
		if strings.Contains(d, "draft: true") {
			continue
		}
		res, err := env.ResolveSymbol(ns, LS(name))
		as.NotNil(res)
		as.Nil(err)
	}
}

func TestBuiltinsHaveDocs(t *testing.T) {
	t.Skip("initially to drive back-filling of documentation")

	as := assert.New(t)
	ns := assert.GetTestEnvironment().GetRoot()

	d := ns.Declared()
	as.NotEqual(0, len(d))

	var missing data.Locals
	for _, name := range d {
		_, err := docstring.Get(string(name))
		if err != nil {
			missing = append(missing, name)
		}
	}

	as.Equal(0, len(missing))
	as.Equal(data.Locals{}, missing)
}

func TestMustGet(t *testing.T) {
	as := assert.New(t)

	d := docstring.MustGet("doc")
	as.NotNil(d)

	defer as.ExpectPanic(
		fmt.Errorf(docstring.ErrSymbolNotDocumented, "blah"),
	)
	_ = docstring.MustGet("blah")
}

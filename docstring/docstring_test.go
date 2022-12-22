package docstring_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/docstring"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestDocString(t *testing.T) {
	as := assert.New(t)

	ifStr, err := docstring.Get("if")
	as.Contains("---", S(ifStr))
	as.Nil(err)

	s, err := docstring.Get("no-way-this-exists")
	as.Empty(s)
	as.EqualError(err, fmt.Sprintf(docstring.ErrDocNotFound, "no-way-this-exists"))
}

func TestDocumentedBuiltinsExist(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()

	for _, name := range docstring.Names() {
		d, _ := docstring.Get(name)
		if strings.Contains(d, "draft: true") {
			continue
		}
		as.NotNil(env.MustResolveSymbol(ns, LS(name)))
	}
}

func TestBuiltinsHaveDocs(t *testing.T) {
	t.Skip("initially to drive back-filling of documentation")

	// starts with ^ or ends (but does not start) with *
	ignorable := regexp.MustCompile(`^(\^.*|[^*]+\*)$`)

	as := assert.New(t)
	ns := assert.GetTestEnvironment().GetRoot()

	d := ns.Declared()
	as.NotEqual(0, len(d))

	var missing data.Names
	for _, name := range d {
		_, err := docstring.Get(string(name))
		if err != nil {
			if ignorable.MatchString(string(name)) {
				continue
			}
			missing = append(missing, name)
		}
	}

	as.Equal(0, len(missing))
	as.Equal(data.Names{}, missing)
}

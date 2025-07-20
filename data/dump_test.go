package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/types"
)

type dumpTest struct{}

func (dumpTest) Name() data.Name {
	return "dump-test"
}

func (dumpTest) Count() int {
	return 42
}

func (dumpTest) Type() ale.Type {
	return types.BasicList
}

func (dumpTest) IsEmpty() bool                           { return false }
func (dumpTest) Car() ale.Value                          { return nil }
func (dumpTest) Cdr() ale.Value                          { return nil }
func (dumpTest) Split() (ale.Value, data.Sequence, bool) { return nil, nil, true }
func (dumpTest) String() string                          { return "" }
func (dumpTest) Equal(ale.Value) bool                    { return false }

var _ interface {
	data.Named
	ale.Typed
	data.Counted
} = dumpTest{}

func TestDumpString(t *testing.T) {
	as := assert.New(t)
	s := S(data.DumpString(dumpTest{}))
	as.Contains(fmt.Sprintf("%s dump-test", data.NameKey), s)
	as.Contains(fmt.Sprintf("%s list", data.TypeKey), s)
	as.Contains(fmt.Sprintf("%s 42", data.CountKey), s)
}

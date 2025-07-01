package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

type dumpTest struct{}

func (dumpTest) Name() data.Local {
	return "dump-test"
}

func (dumpTest) Count() int {
	return 42
}

func (dumpTest) Type() types.Type {
	return types.BasicList
}

func (dumpTest) IsEmpty() bool                            { return false }
func (dumpTest) Car() data.Value                          { return nil }
func (dumpTest) Cdr() data.Value                          { return nil }
func (dumpTest) Split() (data.Value, data.Sequence, bool) { return nil, nil, true }
func (dumpTest) String() string                           { return "" }
func (dumpTest) Equal(data.Value) bool                    { return false }

var _ interface {
	data.Named
	data.Typed
	data.Counted
} = dumpTest{}

func TestDumpString(t *testing.T) {
	as := assert.New(t)
	s := S(data.DumpString(dumpTest{}))
	as.Contains(fmt.Sprintf("%s dump-test", data.NameKey), s)
	as.Contains(fmt.Sprintf("%s list", data.TypeKey), s)
	as.Contains(fmt.Sprintf("%s 42", data.CountKey), s)
}

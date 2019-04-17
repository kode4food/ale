package read_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/read"
)

func TestCreateReader(t *testing.T) {
	as := assert.New(t)
	l := read.Scan("99")
	tr := read.FromScanner(l)
	as.NotNil(tr)
}

func TestReadList(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`(99 "hello" 55.12)`)
	tr := read.FromScanner(l)
	v := tr.First()
	list, ok := v.(*api.List)
	as.True(ok)

	f, r, ok := list.Split()
	as.True(ok)
	as.Integer(99, f)

	f, r, ok = r.Split()
	as.True(ok)
	as.String("hello", f)

	f, r, ok = r.Split()
	as.True(ok)
	as.Float(55.12, f)

	f, r, ok = r.Split()
	as.False(ok)
}

func TestReadVector(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`[99 "hello" 55.12]`)
	tr := read.FromScanner(l)
	v := tr.First()
	vector, ok := v.(api.Vector)
	as.True(ok)

	res, ok := vector.ElementAt(0)
	as.True(ok)
	as.Integer(99, res)

	res, ok = vector.ElementAt(1)
	as.True(ok)
	as.String("hello", res)

	res, ok = vector.ElementAt(2)
	as.True(ok)
	as.Float(55.120, res)
}

func TestReadMap(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`{:name "blah" :age 99}`)
	tr := read.FromScanner(l)
	v := tr.First()
	m, ok := v.(api.Associative)
	as.True(ok)
	as.Integer(2, m.Count())
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`(99 ("hello" "there") 55.12)`)
	tr := read.FromScanner(l)
	v := tr.First()
	list, ok := v.(*api.List)
	as.True(ok)

	f, r, ok := list.Split()
	as.True(ok)
	as.Integer(99, f)

	// get nested list
	f, r, ok = r.Split()
	as.True(ok)
	list2, ok := f.(*api.List)
	as.True(ok)

	// iterate over the rest of top-level list
	f, r, ok = r.Split()
	as.True(ok)
	as.Float(55.12, f)

	f, r, ok = r.Split()
	as.False(ok)

	// iterate over the nested list
	f, r, ok = list2.Split()
	as.True(ok)
	as.String("hello", f)

	f, r, ok = r.Split()
	as.True(ok)
	as.String("there", f)

	f, r, ok = r.Split()
	as.False(ok)
}

func testReaderError(t *testing.T, src string, err error) {
	as := assert.New(t)

	defer as.ExpectPanic(err.Error())

	l := read.Scan(S(src))
	tr := read.FromScanner(l)
	api.Last(tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", fmt.Errorf(read.ListNotClosed))
	testReaderError(t, "[99 100 ", fmt.Errorf(read.VectorNotClosed))
	testReaderError(t, "{:key 99", fmt.Errorf(read.MapNotClosed))

	testReaderError(t, "99 100)", fmt.Errorf(read.UnmatchedListEnd))
	testReaderError(t, "99 100]", fmt.Errorf(read.UnmatchedVectorEnd))
	testReaderError(t, "99}", fmt.Errorf(read.UnmatchedMapEnd))
	testReaderError(t, "{99}", fmt.Errorf(read.MapNotPaired))

	testReaderError(t, "(", fmt.Errorf(read.ListNotClosed))
	testReaderError(t, "'", fmt.Errorf(read.PrefixedNotPaired, "ale/quote"))
	testReaderError(t, "~@", fmt.Errorf(read.PrefixedNotPaired, "ale/unquote-splicing"))
	testReaderError(t, "~", fmt.Errorf(read.PrefixedNotPaired, "ale/unquote"))
}

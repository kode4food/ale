package read_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/read"
	"gitlab.com/kode4food/ale/stdlib"
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

	i := stdlib.Iterate(list)
	val, ok := i.Next()
	as.True(ok)
	as.Integer(99, val)

	val, ok = i.Next()
	as.True(ok)
	as.String("hello", val)

	val, ok = i.Next()
	as.True(ok)
	as.Float(55.12, val)

	_, ok = i.Next()
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

	i1 := stdlib.Iterate(list)
	val, ok := i1.Next()
	as.True(ok)
	as.Integer(99, val)

	// get nested list
	val, ok = i1.Next()
	as.True(ok)
	list2, ok := val.(*api.List)
	as.True(ok)

	// iterate over the rest of top-level list
	val, ok = i1.Next()
	as.True(ok)
	as.Float(55.12, val)

	_, ok = i1.Next()
	as.False(ok)

	// iterate over the nested list
	i2 := stdlib.Iterate(list2)
	val, ok = i2.Next()
	as.True(ok)
	as.String("hello", val)

	val, ok = i2.Next()
	as.True(ok)
	as.String("there", val)

	_, ok = i2.Next()
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

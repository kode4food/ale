package read_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read"
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
	list, ok := v.(data.List)
	as.True(ok)

	f, r, ok := list.Split()
	as.True(ok)
	as.Number(99, f)

	f, r, ok = r.Split()
	as.True(ok)
	as.String("hello", f)

	f, r, ok = r.Split()
	as.True(ok)
	as.Number(55.12, f)

	_, _, ok = r.Split()
	as.False(ok)
}

func TestReadVector(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`[99 "hello" 55.12]`)
	tr := read.FromScanner(l)
	v := tr.First()
	vector, ok := v.(data.Vector)
	as.True(ok)

	res, ok := vector.ElementAt(0)
	as.True(ok)
	as.Number(99, res)

	res, ok = vector.ElementAt(1)
	as.True(ok)
	as.String("hello", res)

	res, ok = vector.ElementAt(2)
	as.True(ok)
	as.Number(55.120, res)
}

func TestReadMap(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`{:name "blah" :age 99}`)
	tr := read.FromScanner(l)
	v := tr.First()
	m, ok := v.(data.Object)
	as.True(ok)
	as.Number(2, m.Count())
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`(99 ("hello" "there") 55.12)`)
	tr := read.FromScanner(l)
	v := tr.First()
	list, ok := v.(data.List)
	as.True(ok)

	f, r, ok := list.Split()
	as.True(ok)
	as.Number(99, f)

	// get nested list
	f, r, ok = r.Split()
	as.True(ok)
	list2, ok := f.(data.List)
	as.True(ok)

	// iterate over the rest of top-level list
	f, r, ok = r.Split()
	as.True(ok)
	as.Number(55.12, f)

	_, r, ok = r.Split()
	as.False(ok)

	// iterate over the nested list
	f, r, ok = list2.Split()
	as.True(ok)
	as.String("hello", f)

	f, r, ok = r.Split()
	as.True(ok)
	as.String("there", f)

	_, _, ok = r.Split()
	as.False(ok)
}

func testReaderError(t *testing.T, src string, err error) {
	as := assert.New(t)

	defer as.ExpectPanic(err.Error())

	l := read.Scan(S(src))
	tr := read.FromScanner(l)
	data.Last(tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", errors.New(read.ErrListNotClosed))
	testReaderError(t, "[99 100 ", errors.New(read.ErrVectorNotClosed))
	testReaderError(t, "{:key 99", errors.New(read.ErrMapNotClosed))

	testReaderError(t, "99 100)", errors.New(read.ErrUnmatchedListEnd))
	testReaderError(t, "99 100]", errors.New(read.ErrUnmatchedVectorEnd))
	testReaderError(t, "99}", errors.New(read.ErrUnmatchedMapEnd))
	testReaderError(t, "{99}", errors.New(data.ErrMapNotPaired))

	testReaderError(t, "(1 2 . 3 4)", errors.New(read.ErrInvalidListSyntax))
	testReaderError(t, "(.)", errors.New(read.ErrInvalidListSyntax))
	testReaderError(t, ".", errors.New(read.ErrUnexpectedDot))

	testReaderError(t, "(", errors.New(read.ErrListNotClosed))
	testReaderError(t, "'", fmt.Errorf(read.ErrPrefixedNotPaired, "ale/quote"))
	testReaderError(t, ",@", fmt.Errorf(read.ErrPrefixedNotPaired, "ale/unquote-splicing"))
	testReaderError(t, ",", fmt.Errorf(read.ErrPrefixedNotPaired, "ale/unquote"))
	testReaderError(t, "~", fmt.Errorf(read.ErrPrefixedNotPaired, "ale/pattern"))
}

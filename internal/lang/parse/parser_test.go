package parse_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/read"
)

func TestReadList(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	tr := read.MustFromString(ns, `(99 "hello" 55.12)`)
	v := tr.Car()
	list, ok := v.(*data.List)
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

func TestReadDotted(t *testing.T) {
	as := assert.New(t)

	ns := assert.GetTestNamespace()
	tr := read.MustFromString(ns, `(99 . 100)`)
	c1 := tr.Car().(*data.Cons)
	as.Number(99, c1.Car())
	as.Number(100, c1.Cdr())

	tr = read.MustFromString(ns, `(99 . (100 101))`)
	l1 := tr.Car().(*data.List)
	as.Number(99, l1.Car())
	as.Number(100, l1.Cdr().(data.Pair).Car())
	as.Number(101, l1.Cdr().(data.Pair).Cdr().(data.Pair).Car())
}

func TestReadVector(t *testing.T) {
	as := assert.New(t)

	ns := assert.GetTestNamespace()
	tr := read.MustFromString(ns, `[99 "hello" 55.12]`)
	v := tr.Car()
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

	ns := assert.GetTestNamespace()
	tr := read.MustFromString(ns, `{:name "blah" :age 99}`)
	v := tr.Car()
	m, ok := v.(*data.Object)
	as.True(ok)
	as.Number(2, m.Count())
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)

	ns := assert.GetTestNamespace()
	tr := read.MustFromString(ns, `(99 ("hello" "there") 55.12)`)
	v := tr.Car()
	list, ok := v.(*data.List)
	as.True(ok)

	f, r, ok := list.Split()
	as.True(ok)
	as.Number(99, f)

	// get nested list
	f, r, ok = r.Split()
	as.True(ok)
	list2, ok := f.(*data.List)
	as.True(ok)

	// iterate over the rest of the top-level list
	f, r, ok = r.Split()
	as.True(ok)
	as.Number(55.12, f)

	_, _, ok = r.Split()
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

func testReaderError(t *testing.T, src, err string, args ...any) {
	as := assert.New(t)

	defer as.ExpectPanic(fmt.Errorf(err, args...))

	ns := assert.GetTestNamespace()
	tr := read.MustFromString(ns, S(src))
	sequence.Last(tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", parse.ErrListNotClosed)
	testReaderError(t, "[99 100 ", parse.ErrVectorNotClosed)
	testReaderError(t, "{:key 99", parse.ErrObjectNotClosed)

	testReaderError(t, "99 100)", parse.ErrUnmatchedListEnd)
	testReaderError(t, "99 100]", parse.ErrUnmatchedVectorEnd)
	testReaderError(t, "99}", parse.ErrUnmatchedObjectEnd)
	testReaderError(t, "{99}", data.ErrMapNotPaired)

	testReaderError(t, "(1 2 . 3 4)", parse.ErrInvalidListSyntax)
	testReaderError(t, "(.)", parse.ErrInvalidListSyntax)
	testReaderError(t, ".", parse.ErrUnexpectedDot)

	testReaderError(t, "(", parse.ErrListNotClosed)
	testReaderError(t, "'", parse.ErrPrefixedNotPaired, "ale/quote")
	testReaderError(t, ",@", parse.ErrPrefixedNotPaired, "ale/unquote-splicing")
	testReaderError(t, ",", parse.ErrPrefixedNotPaired, "ale/unquote")

	testReaderError(t, "//", data.ErrInvalidSymbol, "//")
	testReaderError(t, "/bad", data.ErrInvalidSymbol, "/bad")
	testReaderError(t, "bad/", data.ErrInvalidSymbol, "bad/")
	testReaderError(t, "bad///", data.ErrInvalidSymbol, "bad///")
	testReaderError(t, "ale/er/ror", data.ErrInvalidSymbol, "ale/er/ror")

	testReaderError(t, `"unterminated`, lex.ErrStringNotTerminated)
}

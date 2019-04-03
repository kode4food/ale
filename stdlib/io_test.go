package stdlib_test

import (
	"bytes"
	"io"
	"testing"

	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/stdlib"
)

type mockWriterCloser struct {
	closed bool
	io.Writer
}

func (m *mockWriterCloser) Close() error {
	m.closed = true
	return nil
}

func TestReader(t *testing.T) {
	as := assert.New(t)

	b1 := []byte("12ö34")
	b2 := []byte("12\n34")

	r1 := stdlib.NewReader(bytes.NewReader(b1), stdlib.RuneInput)
	as.True(r1.IsSequence())
	as.String("1", r1.First())
	as.String("2", r1.Rest().First())
	as.String("ö", r1.Rest().Rest().First())
	as.String("3", r1.Rest().Rest().Rest().First())
	as.String("4", r1.Rest().Rest().Rest().Rest().First())
	s1 := r1.Rest().Rest().Rest().Rest().Rest()
	as.False(s1.IsSequence())

	r2 := stdlib.NewReader(bytes.NewReader(b1), stdlib.LineInput)
	as.True(r2.IsSequence())
	as.String("12ö34", r2.First())
	as.False(r2.Rest().IsSequence())

	r3 := stdlib.NewReader(bytes.NewReader(b2), stdlib.LineInput)
	as.True(r3.IsSequence())
	as.String("12", r3.First())
	as.String("34", r3.Rest().First())
	as.False(r3.Rest().Rest().IsSequence())
}

func TestWriter(t *testing.T) {
	as := assert.New(t)

	var buf bytes.Buffer
	c := &mockWriterCloser{
		Writer: &buf,
	}

	w := stdlib.NewWriter(c, stdlib.StrOutput)
	w.Write(S("hello"))
	w.Write(V(S("there"), S("you")))
	w.(stdlib.Closer).Close()

	as.Contains(":type writer", w)
	as.String(`hello["there" "you"]`, buf.String())
	as.True(c.closed)
}

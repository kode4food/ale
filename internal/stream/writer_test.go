package stream_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/stream"
)

type mockWriterCloser struct {
	closed bool
	io.Writer
}

func (m *mockWriterCloser) Close() error {
	m.closed = true
	return nil
}

func TestWriter(t *testing.T) {
	as := assert.New(t)

	var buf bytes.Buffer
	c := &mockWriterCloser{
		Writer: &buf,
	}

	w := stream.NewWriter(c, stream.StrOutput)

	var write data.Lambda
	v, _ := w.Get(stream.WriteKey)
	write = v.(data.Lambda)

	var cl data.Lambda
	v, _ = w.Get(stream.CloseKey)
	cl = v.(data.Lambda)

	write.Call(S("hello"))
	write.Call(V(S("there"), S("you")))
	cl.Call()

	as.String(`hello["there" "you"]`, buf.String())
	as.True(c.closed)
}

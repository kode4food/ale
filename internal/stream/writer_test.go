package stream_test

import (
	"bytes"
	"io"
	"testing"

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
	w.Write(S("hello"))
	w.Write(V(S("there"), S("you")))
	w.(stream.Closer).Close()

	as.Contains(":type writer", w)
	as.String(`hello["there" "you"]`, buf.String())
	as.True(c.closed)
}

package stream_test

import (
	"io"
	"strings"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/data"
)

type mockWriterCloser struct {
	io.Writer
	closed bool
}

func (m *mockWriterCloser) Close() error {
	m.closed = true
	return nil
}

func TestWriter(t *testing.T) {
	as := assert.New(t)

	var buf strings.Builder
	c := &mockWriterCloser{
		Writer: &buf,
	}

	w := stream.NewWriter(c, stream.StrOutput)

	var write data.Procedure
	v, _ := w.Get(stream.WriteKey)
	write = v.(data.Procedure)

	var cl data.Procedure
	v, _ = w.Get(stream.CloseKey)
	cl = v.(data.Procedure)

	write.Call(S("hello"))
	write.Call(V(S("there"), S("you")))
	cl.Call()

	as.String(`hello["there" "you"]`, buf.String())
	as.True(c.closed)
}

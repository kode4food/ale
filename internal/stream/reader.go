package stream

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

// InputFunc is a callback used to unmarshal values from a Reader
type InputFunc func(*bufio.Reader) (data.Value, bool)

const (
	ErrInvalidBlockSize = "block size must be greater than zero: %d"
)

// NewReader wraps a Go Reader, coupling it with an input function
func NewReader(r io.Reader, i InputFunc) data.Sequence {
	var resolver sequence.LazyResolver
	br := bufio.NewReader(r)

	resolver = func() (data.Value, data.Sequence, bool) {
		if v, ok := i(br); ok {
			return v, sequence.NewLazy(resolver), true
		}
		if c, ok := r.(io.Closer); ok {
			_ = c.Close()
		}
		return data.Null, data.Null, false
	}

	return sequence.NewLazy(resolver)
}

// BlockInput creates an InputFunc that reads up to a fixed number of bytes
// from the Reader.
func BlockInput(size int) (InputFunc, error) {
	if size <= 0 {
		return nil, fmt.Errorf(ErrInvalidBlockSize, size)
	}
	return func(r *bufio.Reader) (data.Value, bool) {
		b := make([]byte, size)
		n, err := r.Read(b)
		if err == nil || errors.Is(err, io.EOF) && n > 0 {
			return data.String(b[0:n]), true
		}
		return data.Null, false
	}, nil
}

// LineInput is the standard single line input function
func LineInput(r *bufio.Reader) (data.Value, bool) {
	l, err := r.ReadBytes('\n')
	if err == nil {
		return data.String(l[0 : len(l)-1]), true
	}
	if errors.Is(err, io.EOF) && len(l) > 0 {
		return data.String(l), true
	}
	return data.Null, false
}

// RuneInput is the standard single rune input function
func RuneInput(r *bufio.Reader) (data.Value, bool) {
	if c, _, err := r.ReadRune(); err == nil {
		return data.String(c), true
	}
	return data.Null, false
}

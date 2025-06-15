package asm

import (
	"fmt"

	"github.com/kode4food/ale/pkg/data"
)

// Error messages
const (
	// ErrPairExpected is raised when a vector is provided, but does not
	// contain exactly two elements
	ErrPairExpected = "binding pair expected, got %d elements"
)

const (
	ForEach = data.Local("for-each")

	pairType = "binding pair"
	seqType  = "sequence"
)

func parseForEachCall(
	p *Parser, s data.Sequence,
) (Emit, data.Sequence, error) {
	f, r, ok := s.Split()
	if !ok {
		return nil, nil, typeError(pairType, f)
	}
	k, v, err := parseForEachBinding(f)
	if err != nil {
		return nil, nil, err
	}
	pc := p.withParams(data.Locals{k})
	block, rest, err := pc.block(r)
	if err != nil {
		return nil, nil, err
	}

	return func(e *Encoder) error {
		s, ok := e.resolveEncoderArg(v)
		if !ok {
			return fmt.Errorf(ErrUnexpectedParameter, v)
		}
		seq, err := assertType[data.Sequence](seqType, s)
		if err != nil {
			return err
		}
		for f, r, ok := seq.Split(); ok; f, r, ok = r.Split() {
			if err := block(pc.wrapEncoder(e, f)); err != nil {
				return err
			}
		}
		return nil
	}, rest, nil
}

func parseForEachBinding(v data.Value) (data.Local, data.Value, error) {
	b, err := assertType[data.Vector](pairType, v)
	if err != nil {
		return "", nil, err
	}
	if len(b) != 2 {
		return "", nil, fmt.Errorf(ErrPairExpected, len(b))
	}
	l, err := assertType[data.Local](nameType, b[0])
	if err != nil {
		return "", nil, err
	}
	return l, b[1], nil
}

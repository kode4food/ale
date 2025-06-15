package asm

import (
	"errors"
	"fmt"
	"slices"

	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

type (
	Parser struct {
		calls  namedAsmParsers
		params data.Locals
	}

	namedAsmParsers map[data.Local]asmParse

	asmParse     func(*Parser, data.Sequence) (Emit, data.Sequence, error)
	asmArgsParse func(*Parser, ...data.Value) (Emit, error)
)

const (
	// ErrUnknownDirective is raised when an unknown directive is called
	ErrUnknownDirective = "unknown directive: %s"

	// ErrUnexpectedForm is raised when an unexpected form is encountered in
	// the assembler block
	ErrUnexpectedForm = "unexpected form: %s"

	// ErrTooFewArguments is raised when an instruction is encountered in the
	// assembler block not accompanied by enough operands
	ErrTooFewArguments = "incomplete %s instruction, args expected: %d"

	// ErrExpectedEndOfBlock is raised when an end-of-block marker is expected
	// but an end of stream is encountered instead
	ErrExpectedEndOfBlock = "expected end of block"
)

const (
	EndBlock = data.Local("end")
)

func makeAsmParser(calls namedAsmParsers) *Parser {
	return &Parser{
		calls:  calls,
		params: data.Locals{},
	}
}

func (p *Parser) withParams(n data.Locals) *Parser {
	res := *p
	res.params = slices.Clone(n)
	return &res
}

func (p *Parser) next(s data.Sequence) (Emit, data.Sequence, error) {
	f, r, _ := s.Split()
	switch t := f.(type) {
	case data.Keyword:
		return func(e *Encoder) error {
			e.Emit(isa.Label, e.getLabelIndex(t.Name()))
			return nil
		}, r, nil
	case data.Local:
		parse, ok := p.calls[t]
		if !ok {
			return nil, nil, fmt.Errorf(ErrUnknownDirective, t)
		}
		return parse(p, r)
	default:
		return nil, nil, fmt.Errorf(ErrUnexpectedForm, data.ToString(f))
	}
}

func (p *Parser) sequence(s data.Sequence) (Emit, error) {
	if s.IsEmpty() {
		return noAsmEmit, nil
	}
	next, rest, err := p.next(s)
	if err != nil {
		return nil, err
	}
	return p.rest(next, rest)
}

func (p *Parser) rest(emit Emit, r data.Sequence) (Emit, error) {
	next, err := p.sequence(r)
	if err != nil {
		return nil, err
	}
	return func(e *Encoder) error {
		if err := emit(e); err != nil {
			return err
		}
		return next(e)
	}, nil
}

func (p *Parser) block(s data.Sequence) (Emit, data.Sequence, error) {
	if s.IsEmpty() {
		return nil, nil, errors.New(ErrExpectedEndOfBlock)
	}
	f, r, _ := s.Split()
	if f == EndBlock {
		return noAsmEmit, r, nil
	}
	next, rest, err := p.next(s)
	if err != nil {
		return nil, nil, err
	}
	return p.blockRest(next, rest)
}

func (p *Parser) blockRest(
	emit Emit, r data.Sequence,
) (Emit, data.Sequence, error) {
	next, rest, err := p.block(r)
	if err != nil {
		return nil, nil, err
	}
	return func(e *Encoder) error {
		if err := emit(e); err != nil {
			return err
		}
		return next(e)
	}, rest, nil
}

func parseArgs(inst data.Local, argc int, fn asmArgsParse) asmParse {
	return func(p *Parser, s data.Sequence) (Emit, data.Sequence, error) {
		args, rest, ok := sequence.Take(s, argc)
		if !ok {
			return nil, nil, fmt.Errorf(ErrTooFewArguments, inst, argc)
		}
		res, err := fn(p, args...)
		if err != nil {
			return nil, nil, err
		}
		return res, rest, nil
	}
}

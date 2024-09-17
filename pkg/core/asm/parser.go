package asm

import (
	"errors"
	"fmt"
	"slices"

	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

type (
	asmParser struct {
		calls  namedAsmParsers
		params data.Locals
	}

	namedAsmParsers map[data.Local]asmParse

	asmParse     func(*asmParser, data.Sequence) (asmEmit, data.Sequence, error)
	asmArgsParse func(*asmParser, ...data.Value) (asmEmit, error)
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
	MakeSpecial = data.Local("!make-special")
	EndBlock    = data.Local(".end")
)

func makeAsmParser(calls namedAsmParsers) *asmParser {
	return &asmParser{
		calls:  calls,
		params: data.Locals{},
	}
}

func (p *asmParser) withParams(n data.Locals) *asmParser {
	res := p.copy()
	res.params = slices.Clone(n)
	return res
}

func (p *asmParser) copy() *asmParser {
	res := *p
	res.params = slices.Clone(res.params)
	return &res
}

func (p *asmParser) parse(forms data.Sequence) (asmEmit, error) {
	if f, r, ok := forms.Split(); ok && f == MakeSpecial {
		return p.specialCall(r)
	}
	return p.sequence(forms)
}

func (p *asmParser) next(s data.Sequence) (asmEmit, data.Sequence, error) {
	f, r, _ := s.Split()
	switch t := f.(type) {
	case data.Keyword:
		return func(e *asmEncoder) error {
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

func (p *asmParser) sequence(s data.Sequence) (asmEmit, error) {
	if s.IsEmpty() {
		return noAsmEmit, nil
	}
	next, rest, err := p.next(s)
	if err != nil {
		return nil, err
	}
	return p.rest(next, rest)
}

func (p *asmParser) rest(emit asmEmit, r data.Sequence) (asmEmit, error) {
	next, err := p.sequence(r)
	if err != nil {
		return nil, err
	}
	return func(e *asmEncoder) error {
		if err := emit(e); err != nil {
			return err
		}
		return next(e)
	}, nil
}

func (p *asmParser) block(s data.Sequence) (asmEmit, data.Sequence, error) {
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

func (p *asmParser) blockRest(
	emit asmEmit, r data.Sequence,
) (asmEmit, data.Sequence, error) {
	next, rest, err := p.block(r)
	if err != nil {
		return nil, nil, err
	}
	return func(e *asmEncoder) error {
		if err := emit(e); err != nil {
			return err
		}
		return next(e)
	}, rest, nil
}

func parseArgs(inst data.Local, argc int, fn asmArgsParse) asmParse {
	return func(p *asmParser, s data.Sequence) (asmEmit, data.Sequence, error) {
		args, rest, ok := take(s, argc)
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

func take(s data.Sequence, count int) (data.Vector, data.Sequence, bool) {
	var f data.Value
	var ok bool
	res := make(data.Vector, count)
	for i := 0; i < count; i++ {
		if f, s, ok = s.Split(); !ok {
			return nil, nil, false
		}
		res[i] = f
	}
	return res, s, true
}
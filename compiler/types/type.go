package types

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/data"
)

type (
	// Type is the basic interface for type checking
	Type interface {
		Name() data.Name
		Satisfies(Type) error
	}

	// Boolean represents the boolean native type
	Boolean struct{}

	// Numeric represents the numeric native type
	Numeric struct{}

	// String represents the string native type
	String struct{}

	// Null represents the empty list or nil value as a distinct type
	Null struct{}

	// Composite represents a composite type (list, tuple, record, sum)
	Composite interface {
		Type
		Composite()
	}

	// List represents a list type, where each element is the same type
	List struct {
		element Type
	}

	// Tuple represents a fixed set of independently typed values
	Tuple struct {
		elements []Type
	}

	// RecordEntry represents a record's named fields
	RecordEntry struct {
		Name data.Name
		Type Type
	}

	// Record represents a data structure with named fields
	Record []RecordEntry

	// Sum represents a union of types
	Sum struct {
		types []Type
	}

	// LambdaCase represents a lambda case signature
	LambdaCase struct {
		args   []Type
		rest   bool
		result Type
	}

	// Lambda represents a function or closure
	Lambda struct {
		cases []LambdaCase
	}
)

func (b *Boolean) Name() data.Name {
	return "bool"
}

func (b *Boolean) Satisfies(t Type) error {
	if _, ok := t.(*Boolean); ok {
		return nil
	}
	return errors.New("type does not satisfy bool")
}

func (n *Numeric) Name() data.Name {
	return "number"
}

func (n *Numeric) Satisfies(t Type) error {
	if _, ok := t.(*Numeric); ok {
		return nil
	}
	return errors.New("type does not satisfy number")
}

func (s *String) Name() data.Name {
	return "string"
}

func (s *String) Satisfies(t Type) error {
	if _, ok := t.(*String); ok {
		return nil
	}
	return errors.New("type does not satisfy string")
}

func (n *Null) Name() data.Name {
	return "null"
}

func (n *Null) Satisfies(t Type) error {
	if _, ok := t.(*Null); ok {
		return nil
	}
	return errors.New("type does not satisfy null")
}

func (l *List) Name() data.Name {
	sName := string(l.element.Name())
	var res = data.Name(fmt.Sprintf("list<%s>", sName))
	return res
}

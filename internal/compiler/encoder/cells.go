package encoder

import (
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

type (
	// CellType marks a cell as having a certain disposition
	CellType int

	// Cell attaches a name to a type/disposition
	Cell struct {
		Name data.Local
		Type CellType
	}

	// IndexedCells encapsulates a group of IndexedCells
	IndexedCells []*IndexedCell

	// IndexedCell attaches an Index to a Cell
	IndexedCell struct {
		*Cell
		Index isa.Operand
	}

	// ScopedCell attaches a Scope to a Cell
	ScopedCell struct {
		*Cell
		Scope
	}
)

// Cell dispositions
const (
	ValueCell CellType = iota
	ReferenceCell
	RestCell
)

func newCell(t CellType, n data.Local) *Cell {
	return &Cell{
		Name: n,
		Type: t,
	}
}

func newIndexedCell(i isa.Operand, c *Cell) *IndexedCell {
	return &IndexedCell{
		Cell:  c,
		Index: i,
	}
}

func newScopedCell(s Scope, c *Cell) *ScopedCell {
	return &ScopedCell{
		Scope: s,
		Cell:  c,
	}
}

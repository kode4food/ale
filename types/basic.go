package types

type (
	// BasicType accepts a value of its own Type
	BasicType interface {
		Type
		basic() // marker
	}

	basic struct {
		name string
	}
)

// Basic returns a Type that accepts a Value of its own Type
func Basic(name string) BasicType {
	return &basic{
		name: name,
	}
}

func (*basic) basic() {}

func (b *basic) Name() string {
	return b.name
}

func (b *basic) Accepts(other Type) bool {
	return b == other
}

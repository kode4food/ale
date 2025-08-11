//go:generate go run github.com/kode4food/gen-maxkind
package ffi

import (
	"errors"
	"reflect"
	"sync"
	"unsafe"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi/maxkind"
)

type (
	// Wrapper can marshal a native Go value to and from a data.Value
	Wrapper interface {
		// Wrap converts a native Go value to an ale.Value
		Wrap(*Context, reflect.Value) (ale.Value, error)

		// Unwrap converts an ale.Value back to a native Go value
		Unwrap(ale.Value) (reflect.Value, error)
	}

	typeCache struct {
		entries map[reflect.Type]Wrapper
		sync.RWMutex
	}

	handler func(reflect.Type) (Wrapper, error)
)

// ErrUnsupportedType is raised when wrapping encounters an unsupported type
const ErrUnsupportedType = "unsupported type"

var (
	cache = makeTypeCache()

	_zero = reflect.Value{}

	handlers     [maxkind.Value + 1]handler
	handlersOnce sync.Once
)

// Wrap takes a native Go value, potentially builds a Wrapper for its type, and
// returns a marshaled data.Value from the Wrapper
func Wrap(i any) (ale.Value, error) {
	v := reflect.ValueOf(i)
	w, err := WrapType(v.Type())
	if err != nil {
		return data.Null, err
	}
	return w.Wrap(new(Context), v)
}

// MustWrap wraps a Go value into a data.Value or explodes violently
func MustWrap(i any) ale.Value {
	res, err := Wrap(i)
	if err != nil {
		panic(err)
	}
	return res
}

// WrapType potentially builds a Wrapper for the provided reflected Type
func WrapType(t reflect.Type) (Wrapper, error) {
	if w, ok := cache.get(t); ok {
		return w, nil
	}

	// register a stub to avoid wrap cycles
	s := new(struct{ Wrapper })
	cache.put(t, s)

	// register the final Wrapper, and wire it into the stub for those Wrappers
	// that may refer to it
	w, err := makeWrappedType(t)
	if err != nil {
		return nil, err
	}
	cache.put(t, w)
	s.Wrapper = w
	return w, nil
}

func makeWrappedType(t reflect.Type) (Wrapper, error) {
	if t.Implements(dataValue) {
		return wrapDataValue(t)
	}

	handlers := getKindHandlers()
	if handler := handlers[t.Kind()]; handler != nil {
		return handler(t)
	}

	return nil, errors.New(ErrUnsupportedType)
}

func makeTypeCache() *typeCache {
	return &typeCache{
		entries: map[reflect.Type]Wrapper{},
	}
}

func (c *typeCache) get(t reflect.Type) (Wrapper, bool) {
	c.RLock()
	defer c.RUnlock()
	w, ok := c.entries[t]
	return w, ok
}

func (c *typeCache) put(t reflect.Type, w Wrapper) {
	c.Lock()
	defer c.Unlock()
	c.entries[t] = w
}

func getKindHandlers() *[maxkind.Value + 1]handler {
	handlersOnce.Do(func() {
		handlers[reflect.Bool] = makeHandler(boolWrapper{})
		handlers[reflect.Int] = makeHandler(intWrapper[int]{})
		handlers[reflect.Int8] = makeHandler(intWrapper[int8]{})
		handlers[reflect.Int16] = makeHandler(intWrapper[int16]{})
		handlers[reflect.Int32] = makeHandler(intWrapper[int32]{})
		handlers[reflect.Int64] = makeHandler(intWrapper[int64]{})
		handlers[reflect.Uint] = makeHandler(uintWrapper[uint]{})
		handlers[reflect.Uint8] = makeHandler(uintWrapper[uint8]{})
		handlers[reflect.Uint16] = makeHandler(uintWrapper[uint16]{})
		handlers[reflect.Uint32] = makeHandler(uintWrapper[uint32]{})
		handlers[reflect.Uint64] = makeHandler(uintWrapper[uint64]{})
		handlers[reflect.Uintptr] = makeHandler(boxedWrapper[uintptr]{})
		handlers[reflect.Float32] = makeHandler(floatWrapper[float32]{})
		handlers[reflect.Float64] = makeHandler(floatWrapper[float64]{})
		handlers[reflect.Complex64] = makeHandler(complexWrapper[complex64]{})
		handlers[reflect.Complex128] = makeHandler(complexWrapper[complex128]{})
		handlers[reflect.Array] = makeWrappedArray
		handlers[reflect.Chan] = makeWrappedChannel
		handlers[reflect.Func] = makeWrappedFunc
		handlers[reflect.Interface] = makeWrappedInterface
		handlers[reflect.Map] = makeWrappedMap
		handlers[reflect.Ptr] = makeWrappedPointer
		handlers[reflect.Slice] = makeWrappedSlice
		handlers[reflect.String] = makeHandler(stringWrapper{})
		handlers[reflect.Struct] = makeWrappedStruct
		handlers[reflect.UnsafePointer] = makeHandler(
			boxedWrapper[unsafe.Pointer]{},
		)
	})
	return &handlers
}

func makeHandler(w Wrapper) handler {
	return func(reflect.Type) (Wrapper, error) {
		return w, nil
	}
}

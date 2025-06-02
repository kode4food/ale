package ffi

import (
	"errors"
	"reflect"
	"sync"
	"unsafe"

	"github.com/kode4food/ale/pkg/data"
)

type (
	// Wrapper can marshal a native Go value to and from a data.Value
	Wrapper interface {
		Wrap(*Context, reflect.Value) (data.Value, error)
		Unwrap(data.Value) (reflect.Value, error)
	}

	typeCache struct {
		entries map[reflect.Type]Wrapper
		sync.RWMutex
	}
)

// ErrUnsupportedType is raised when wrapping encounters an unsupported type
const ErrUnsupportedType = "unsupported type"

var (
	cache = makeTypeCache()

	_emptyValue = reflect.Value{}
)

// Wrap takes a native Go value, potentially builds a Wrapper for its type, and
// returns a marshaled data.Value from the Wrapper
func Wrap(i any) (data.Value, error) {
	v := reflect.ValueOf(i)
	w, err := WrapType(v.Type())
	if err != nil {
		return data.Null, err
	}
	return w.Wrap(new(Context), v)
}

// MustWrap wraps a Go value into a data.Value or explodes violently
func MustWrap(i any) data.Value {
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
	switch t.Kind() {
	case reflect.Bool:
		return boolWrapper{}, nil
	case reflect.Int:
		return intWrapper[int]{}, nil
	case reflect.Int8:
		return intWrapper[int8]{}, nil
	case reflect.Int16:
		return intWrapper[int16]{}, nil
	case reflect.Int32:
		return intWrapper[int32]{}, nil
	case reflect.Int64:
		return intWrapper[int64]{}, nil
	case reflect.Uint:
		return uint64Wrapper[uint]{}, nil
	case reflect.Uint8:
		return uintWrapper[uint8]{}, nil
	case reflect.Uint16:
		return uintWrapper[uint16]{}, nil
	case reflect.Uint32:
		return uintWrapper[uint32]{}, nil
	case reflect.Uint64:
		return uint64Wrapper[uint64]{}, nil
	case reflect.Uintptr:
		return boxedWrapper[uintptr]{}, nil
	case reflect.Float32:
		return floatWrapper[float32]{}, nil
	case reflect.Float64:
		return floatWrapper[float64]{}, nil
	case reflect.Complex64:
		return complexWrapper[complex64]{}, nil
	case reflect.Complex128:
		return complexWrapper[complex128]{}, nil
	case reflect.Array:
		return makeWrappedArray(t)
	case reflect.Chan:
		return makeWrappedChannel(t)
	case reflect.Func:
		return makeWrappedFunc(t)
	case reflect.Interface:
		return makeWrappedInterface(t)
	case reflect.Map:
		return makeWrappedMap(t)
	case reflect.Ptr:
		return makeWrappedPointer(t)
	case reflect.Slice:
		return makeWrappedSlice(t)
	case reflect.String:
		return stringWrapper{}, nil
	case reflect.Struct:
		return makeWrappedStruct(t)
	case reflect.UnsafePointer:
		return boxedWrapper[unsafe.Pointer]{}, nil
	default:
		return nil, errors.New(ErrUnsupportedType)
	}
}

func makeTypeCache() *typeCache {
	return &typeCache{
		entries: map[reflect.Type]Wrapper{},
	}
}

func (c *typeCache) get(t reflect.Type) (Wrapper, bool) {
	c.RLock()
	w, ok := c.entries[t]
	c.RUnlock()
	return w, ok
}

func (c *typeCache) put(t reflect.Type, w Wrapper) {
	c.Lock()
	c.entries[t] = w
	c.Unlock()
}

func zero[T any]() reflect.Value {
	var z T
	return reflect.ValueOf(z)
}

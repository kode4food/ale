package ffi

import (
	"errors"
	"reflect"
	"sync"

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

/*
Unsupported Kinds:
  - Uintptr
  - UnsafePointer
*/
func makeWrappedType(t reflect.Type) (Wrapper, error) {
	if t.Implements(dataValue) {
		return wrapDataValue(t)
	}
	switch t.Kind() {
	case reflect.Array:
		return makeWrappedArray(t)
	case reflect.Slice:
		return makeWrappedSlice(t)
	case reflect.Bool:
		return makeWrappedBool(t), nil
	case reflect.Chan:
		return makeWrappedChannel(t)
	case reflect.Complex64, reflect.Complex128:
		return makeWrappedComplex(t), nil
	case reflect.Float32, reflect.Float64:
		return makeWrappedFloat(t), nil
	case reflect.Func:
		return makeWrappedFunc(t)
	case reflect.Interface:
		return makeWrappedInterface(t)
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return makeWrappedInt(t), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return makeWrappedUnsignedInt(t), nil
	case reflect.Map:
		return makeWrappedMap(t)
	case reflect.Ptr:
		return makeWrappedPointer(t)
	case reflect.String:
		return makeWrappedString(t), nil
	case reflect.Struct:
		return makeWrappedStruct(t)
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

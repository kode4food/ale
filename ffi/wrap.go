package ffi

import (
	"reflect"
	"sync"

	"github.com/kode4food/ale/data"
)

type (
	// Wrapper can marshal a native Go value to and from a data.Value
	Wrapper interface {
		Wrap(*WrapContext, reflect.Value) data.Value
		Unwrap(*UnwrapContext, data.Value) reflect.Value
	}

	typeCache struct {
		sync.RWMutex
		entries map[reflect.Type]Wrapper
	}
)

var cache = makeTypeCache()

// Wrap takes a native Go value, potentially builds a Wrapper for
// its type and then returns a marshalled data.Value from the Wrapper
func Wrap(i interface{}) data.Value {
	if i == nil {
		return data.Nil
	}
	if d, ok := i.(data.Value); ok {
		return d
	}
	v := reflect.ValueOf(i)
	w := wrapType(v.Type())
	c := &WrapContext{}
	return w.Wrap(c, v)
}

func wrapType(t reflect.Type) Wrapper {
	if w, ok := cache.get(t); ok {
		return w
	}

	// register a stub to avoid wrap cycles
	s := &struct{ Wrapper }{}
	cache.put(t, s)

	// register the final Wrapper, and wire it into the
	// stub for those Wrappers that may refer to it
	w := makeWrappedType(t)
	cache.put(t, w)
	s.Wrapper = w
	return w
}

/*
	Unsupported Kinds:
	  * Uintptr
	  * Chan
	  * Complex64
	  * Complex128
	  * UnsafePointer
*/
func makeWrappedType(t reflect.Type) Wrapper {
	switch t.Kind() {
	case reflect.Array:
		return makeWrappedArray(t)
	case reflect.Slice:
		return makeWrappedSlice(t)
	case reflect.Bool:
		return makeWrappedBool(t)
	case reflect.Float32, reflect.Float64:
		return makeWrappedFloat(t)
	case reflect.Func:
		return makeWrappedFunc(t)
	case reflect.Interface:
		return makeWrappedInterface(t)
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8,
		reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32,
		reflect.Int64, reflect.Uint64:
		return makeWrappedInt(t)
	case reflect.Map:
		return makeWrappedMap(t)
	case reflect.Ptr:
		return makeWrappedPointer(t)
	case reflect.String:
		return makeWrappedString(t)
	case reflect.Struct:
		return makeWrappedStruct(t)
	default:
		panic("unsupported type")
	}
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

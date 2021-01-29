package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	// WrapContext tracks items that have already been wrapped
	WrapContext struct {
		mappings []*wrapMapping
	}

	wrapMapping struct {
		From reflect.Value
		To   data.Value
	}

	// UnwrapContext tracks items that have already been unwrapped
	UnwrapContext struct {
		mappings []*unwrapMapping
	}

	unwrapMapping struct {
		From data.Value
		To   reflect.Value
	}
)

var emptyReflectValue = reflect.Value{}

// Get returns a mapped item from the WrapContext list
func (v *WrapContext) Get(from reflect.Value) (data.Value, bool) {
	if !from.IsValid() {
		return data.Nil, true
	}
	for i := len(v.mappings) - 1; i >= 0; i-- {
		m := v.mappings[i]
		if reflect.DeepEqual(m.From, from) {
			return m.To, true
		}
	}
	return nil, false
}

// Put stores a mapping in the WrapContext list
func (v *WrapContext) Put(from reflect.Value, to data.Value) {
	v.mappings = append(v.mappings, &wrapMapping{
		From: from,
		To:   to,
	})
}

// Get returns a mapped item from the WrapContext list
func (v *UnwrapContext) Get(from data.Value) (reflect.Value, bool) {
	for i := len(v.mappings) - 1; i >= 0; i-- {
		m := v.mappings[i]
		if m.From == from {
			return m.To, true
		}
	}
	return emptyReflectValue, false
}

// Put stores a mapping in the WrapContext list
func (v *UnwrapContext) Put(from data.Value, to reflect.Value) {
	v.mappings = append(v.mappings, &unwrapMapping{
		From: from,
		To:   to,
	})
}

package ffi

import (
	"errors"
	"reflect"
)

// Context tracks wrapped reflect.Values to detect cycles
type Context struct {
	parent *Context
	child  bool
	value  reflect.Value
}

// Error messages
const (
	ErrCycleDetected = "cycle detected in wrapping"
)

var _emptyValue = reflect.Value{}

// Push creates a new Context, checking the parent chain for cycles
func (c *Context) Push(v reflect.Value) (*Context, error) {
	if err := c.checkDuplicate(v); err != nil {
		return nil, err
	}
	return &Context{
		parent: c,
		child:  true,
		value:  v,
	}, nil
}

func (c *Context) checkDuplicate(v reflect.Value) error {
	if !c.child {
		return nil
	}
	cv := c.value
	if cv.IsValid() && v.IsValid() && cv.Type() == v.Type() {
		switch cv.Kind() {
		case reflect.Ptr:
			if cv.Pointer() == v.Pointer() {
				return errors.New(ErrCycleDetected)
			}
		case reflect.Slice, reflect.Map:
			if cv.IsNil() || v.IsNil() {
				break
			}
			if cv.Len() == v.Len() && cv.Pointer() == v.Pointer() {
				return errors.New(ErrCycleDetected)
			}
		}
	}
	return c.parent.checkDuplicate(v)
}

package ffi_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/ffi"
)

func TestByteArrayWrapper(t *testing.T) {
	as := assert.New(t)

	id1 := uuid.New()
	s, err := ffi.Wrap(id1)
	as.Nil(err)
	as.Equal(id1.String(), string(s.(data.String)))

	w, err := ffi.WrapType(reflect.TypeOf(id1))
	as.NotNil(w)
	as.Nil(err)

	id2, err := w.Unwrap(s.(data.String))
	as.Nil(err)
	as.Equal(id1, id2.Interface())
}

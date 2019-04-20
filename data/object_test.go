package data_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestObject(t *testing.T) {
	as := assert.New(t)

	o1 := data.Object{
		K("parent"): S("i am the parent"),
		K("name"):   S("parent"),
	}

	o2 := o1.Extend(data.Object{
		K("child"): S("i am the child"),
		K("name"):  S("child"),
	})

	as.String("i am the parent", o2.MustGet(K("parent")))
	as.String("child", o2.MustGet(K("name")))
	as.String("parent", o1.MustGet(K("name")))

	as.Contains(`:name "child"`, o2)
	as.Contains(`:child "i am the child"`, o2)
	as.Contains(`:parent "i am the parent"`, o2)

	defer as.ExpectPanic(fmt.Sprintf(data.ValueNotFound, ":missing"))
	o2.MustGet(K("missing"))
}

func TestObjectCaller(t *testing.T) {
	as := assert.New(t)

	o1 := data.Object{
		K("parent"): S("i am the parent"),
		K("name"):   S("parent"),
	}

	c1 := o1.Caller()
	as.String("i am the parent", c1(K("parent")))
	as.Nil(c1(K("missing")))
	as.String("defaulted", c1(K("missing"), S("defaulted")))
}

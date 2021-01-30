package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestObject(t *testing.T) {
	as := assert.New(t)

	o1 := data.Object{
		K("parent"): S("i am the parent"),
		K("name"):   S("parent"),
	}

	o2 := o1.Merge(data.Object{
		K("child"): S("i am the child"),
		K("name"):  S("child"),
	})

	as.String("i am the parent", o2.MustGet(K("parent")))
	as.String("child", o2.MustGet(K("name")))
	as.String("parent", o1.MustGet(K("name")))

	as.Contains(`:name "child"`, o2)
	as.Contains(`:child "i am the child"`, o2)
	as.Contains(`:parent "i am the parent"`, o2)

	defer as.ExpectPanic(fmt.Sprintf(data.ErrValueNotFound, ":missing"))
	o2.MustGet(K("missing"))
}

func TestObjectCaller(t *testing.T) {
	as := assert.New(t)

	o1 := data.Object{
		K("parent"): S("i am the parent"),
		K("name"):   S("parent"),
	}

	as.String("i am the parent", o1.Call(K("parent")))
	as.Nil(o1.Call(K("missing")))
	as.String("defaulted", o1.Call(K("missing"), S("defaulted")))
}

func TestObjectIterate(t *testing.T) {
	as := assert.New(t)

	o1 := data.Object{
		K("second"): S("second value"),
		K("first"):  S("first value"),
	}
	as.Equal(2, len(o1))

	f1, r1, ok := o1.Split()
	as.True(ok)
	as.Equal(K("first"), f1.(data.Cons).Car())
	as.Equal(S("first value"), f1.(data.Cons).Cdr())
	as.Equal(1, len(r1.(data.Object)))

	f2, r2, ok := r1.Split()
	as.True(ok)
	as.Equal(K("second"), f2.(data.Cons).Car())
	as.Equal(S("second value"), f2.(data.Cons).Cdr())

	_, _, ok = r2.Split()
	as.False(ok)
}

func TestObjectSplitDeterminism(t *testing.T) {
	as := assert.New(t)
	o := data.Object{
		K("z"): I(1024),
		K("x"): I(5),
		K("y"): I(99),
	}
	f1, r1, ok := o.Split()
	r1Str := r1.String()
	as.True(ok)
	for i := 0; i < 50; i++ {
		f2, r2, ok := o.Split()
		as.True(ok)
		as.Equal(f1, f2)
		as.Equal(r1Str, r2.String())
	}
}

func TestObjectEquality(t *testing.T) {
	as := assert.New(t)
	o1 := data.Object{
		K("z"): I(1024),
		K("x"): I(5),
		K("y"): I(99),
	}
	o2 := data.Object{ // Content same
		K("z"): I(1024),
		K("x"): I(5),
		K("y"): I(99),
	}
	o3 := data.Object{ // Missing key
		K("z"): I(1024),
		K("y"): I(99),
	}
	o4 := data.Object{ // Additional Key
		K("z"): I(1024),
		K("x"): I(5),
		K("y"): I(99),
		K("g"): I(1024),
	}
	o5 := data.Object{ // Modified Value in x
		K("z"): I(1024),
		K("x"): I(6),
		K("y"): I(99),
	}
	as.True(o1.Equal(o1))
	as.True(o1.Equal(o2))
	as.False(o1.Equal(o3))
	as.False(o1.Equal(o4))
	as.False(o1.Equal(o5))
	as.False(o1.Equal(I(32)))
}

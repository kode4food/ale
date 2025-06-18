package stream_test

import (
	"embed"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/data"
)

//go:embed test_assets/*
var assets embed.FS

func TestWrapFileSystemList(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	as.NotNil(fs)

	c, ok := fs.Get(stream.ListKey)
	as.True(ok)
	res := c.(data.Procedure).Call(S("test_assets"))
	as.NotNil(res)
	as.Equal(O(
		C(S("test1.txt"), stream.File),
		C(S("test2.txt"), stream.File),
		C(S("dir1"), stream.Dir),
		C(S("dir2"), stream.Dir),
	), res)
}

func TestWrapFileSystemReadAll(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	as.NotNil(fs)

	c, ok := fs.Get(stream.OpenKey)
	as.True(ok)

	res := c.(data.Procedure).Call(S("test_assets/test1.txt"), stream.ReadAll)
	b, ok := res.(data.Bytes)
	as.True(ok)
	as.String("test file 1 in root\n", data.String(b))
}

func TestWrapFileSystemReadLines(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	as.NotNil(fs)

	c, ok := fs.Get(stream.OpenKey)
	as.True(ok)

	r := c.(data.Procedure).Call(S("test_assets/test2.txt"), stream.ReadLines)
	res := sequence.ToVector(r.(data.Sequence))
	as.Equal(2, res.Count())
	as.String("test file 2 in root", res[0])
	as.String("has multiple lines", res[1])
}

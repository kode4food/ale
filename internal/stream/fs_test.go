package stream_test

import (
	"embed"
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
)

//go:embed test_assets/*
var assets embed.FS

func TestWrapFileSystemList(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	if as.NotNil(fs) {
		c, ok := fs.Get(stream.ListKey)
		as.True(ok)
		res := c.(data.Procedure).Call(S("test_assets"))
		if as.NotNil(res) {
			as.Equal(O(
				C(S("test1.txt"), stream.File),
				C(S("test2.txt"), stream.File),
				C(S("dir1"), stream.Dir),
				C(S("dir2"), stream.Dir),
			), res)
		}

		as.Panics(
			func() { _ = c.(data.Procedure).Call(S("not-a-dir")) },
			fmt.Errorf(stream.ErrExpectedDirectory, "not-a-dir"),
		)
	}
}

func TestWrapFileSystemReadAll(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	if as.NotNil(fs) {
		c, ok := fs.Get(stream.OpenKey)
		as.True(ok)

		res := c.(data.Procedure).Call(
			S("test_assets/test1.txt"), stream.ReadAll,
		)
		b, ok := res.(data.Bytes)
		as.True(ok)
		as.String("test file 1 in root\n", data.String(b))
	}
}

func TestWrapFileSystemReadLines(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	if as.NotNil(fs) {
		c, ok := fs.Get(stream.OpenKey)
		as.True(ok)

		r := c.(data.Procedure).Call(
			S("test_assets/test2.txt"), stream.ReadLines,
		)
		res := sequence.ToVector(r.(data.Sequence))
		as.Equal(2, res.Count())
		as.String("test file 2 in root", res[0])
		as.String("has multiple lines", res[1])
	}
}

func TestWrapFileSystemReadStringAndRunes(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	if as.NotNil(fs) {
		c, ok := fs.Get(stream.OpenKey)
		as.True(ok)

		r := c.(data.Procedure).Call(S("test_assets/test1.txt"))
		s := sequence.ToString(r.(data.Sequence))
		as.String("test file 1 in root\n", s)

		v := c.(data.Procedure).Call(
			S("test_assets/test1.txt"), stream.ReadString,
		)
		as.String("test file 1 in root\n", v)
	}
}

func TestWrapFileSystemReadBlocks(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	if as.NotNil(fs) {
		c, ok := fs.Get(stream.OpenKey)
		as.True(ok)

		r1 := c.(data.Procedure).Call(
			S("test_assets/test1.txt"), stream.ReadBlocks,
		)
		v1 := sequence.ToVector(r1.(data.Sequence))
		as.Equal(1, v1.Count())
		as.String("test file 1 in root\n", data.String(v1[0].(data.Bytes)))

		r2 := c.(data.Procedure).Call(
			S("test_assets/test1.txt"), stream.ReadBlocks, I(4),
		)
		v2 := sequence.ToVector(r2.(data.Sequence))
		as.Equal(5, v2.Count())
		as.String("test", data.String(v2[0].(data.Bytes)))
		as.String(" fil", data.String(v2[1].(data.Bytes)))
		as.String("e 1 ", data.String(v2[2].(data.Bytes)))
		as.String("in r", data.String(v2[3].(data.Bytes)))
		as.String("oot\n", data.String(v2[4].(data.Bytes)))
	}
}

func TestWrapFileSystemOpenErrors(t *testing.T) {
	as := assert.New(t)
	fs := stream.WrapFileSystem(assets)
	if as.NotNil(fs) {
		c, ok := fs.Get(stream.OpenKey)
		as.True(ok)
		open := c.(data.Procedure)

		as.Panics(func() {
			_ = open.Call(S("test_assets/test1.txt"), K("unknown-mode"))
		})
		as.Panics(func() {
			_ = open.Call(S("test_assets/test1.txt"), stream.ReadBlocks, I(0))
		})
		as.Panics(func() {
			_ = open.Call(
				S("test_assets/test1.txt"), stream.ReadBlocks, I(4), I(8),
			)
		})
		as.Panics(func() {
			_ = open.Call(S("test_assets/dir1"), stream.ReadAll)
		})
		as.Panics(func() {
			_ = open.Call(S("test_assets/no-such-file.txt"), stream.ReadAll)
		})
	}
}

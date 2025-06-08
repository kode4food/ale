package stream

import (
	"fmt"
	"io/fs"

	"github.com/kode4food/ale/pkg/data"
)

const (
	ListKey = data.Keyword("list")
	Dir     = data.Keyword("dir")
	File    = data.Keyword("file")

	OpenKey    = data.Keyword("open")
	ReadLines  = data.Keyword("read-lines")
	ReadAll    = data.Keyword("read-all")
	ReadBlocks = data.Keyword("read-blocks")
)

const (
	ErrExpectedDirectory    = "expected a directory, got a file: %s"
	ErrUnknownOpenMode      = "unknown open mode: %s"
	ErrUnexpectedReadLength = "expected to read %d bytes, got %d"
	ErrExpectedFile         = "expected a file, got a directory: %s"
	ErrExpectedBlockSize    = "expected a block size, got: %s"
)

func WrapFileSystem(fileSystem fs.FS) *data.Object {
	return data.NewObject(
		data.NewCons(ListKey, bindList(fileSystem)),
		data.NewCons(OpenKey, bindOpen(fileSystem)),
	)
}

func bindList(fileSystem fs.FS) data.Call {
	return func(args ...data.Value) data.Value {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			panic(err)
		}
		path := args[0].(data.String)
		f, err := fileSystem.Open(path.String())
		if err != nil {
			panic(err)
		}
		defer func() { _ = f.Close() }()
		s, err := f.Stat()
		if err == nil || s.IsDir() {
			res, err := listDirectory(f)
			if err != nil {
				panic(err)
			}
			return res
		}
		panic(fmt.Errorf(ErrExpectedDirectory, path))
	}
}

func listDirectory(f fs.File) (*data.Object, error) {
	entries, err := f.(fs.ReadDirFile).ReadDir(-1)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return data.EmptyObject, nil
	}
	res := data.EmptyObject
	for _, e := range entries {
		res = res.Put(
			data.NewCons(data.String(e.Name()), getDirEntryType(e)),
		).(*data.Object)
	}
	return res, nil
}

func getDirEntryType(s fs.DirEntry) data.Keyword {
	if s.IsDir() {
		return Dir
	}
	return File
}

func bindOpen(fs fs.FS) data.Call {
	return func(args ...data.Value) data.Value {
		if err := data.CheckRangedArity(1, 3, len(args)); err != nil {
			panic(err)
		}
		path := args[0].(data.String)
		f, s, err := openFile(fs, path.String())
		if err != nil {
			panic(err)
		}
		return createReader(f, s, args[1:]...)
	}
}

func createReader(f fs.File, s fs.FileInfo, args ...data.Value) data.Value {
	if len(args) == 0 {
		return NewReader(f, RuneInput)
	}
	switch args[0] {
	case ReadLines:
		return NewReader(f, LineInput)
	case ReadAll:
		return readAll(f, s.Size())
	case ReadBlocks:
		if len(args) < 2 {
			panic(fmt.Errorf(ErrExpectedBlockSize, data.Null))
		}
		size := args[1].(data.Integer)
		input, err := BlockInput(int(size))
		if err != nil {
			panic(err)
		}
		return NewReader(f, input)
	default:
		panic(fmt.Errorf(ErrUnknownOpenMode, args[0]))
	}
}

func openFile(fs fs.FS, path string) (fs.File, fs.FileInfo, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, nil, err
	}
	s, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}
	if s.IsDir() {
		_ = f.Close()
		return nil, nil, fmt.Errorf(ErrExpectedFile, path)
	}
	return f, s, nil
}

func readAll(f fs.File, size int64) data.Value {
	buf := make([]byte, size)
	l, err := f.Read(buf)
	if err != nil {
		panic(err)
	}
	if int64(l) != size {
		panic(fmt.Errorf(ErrUnexpectedReadLength, size, l))
	}
	_ = f.Close()
	return data.String(buf[:l])
}

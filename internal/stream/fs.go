package stream

import (
	"fmt"
	"io/fs"
	"unsafe"

	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

type FileSystem struct {
	fs   fs.FS
	List data.Procedure
	Open data.Procedure
}

const (
	ListKey = data.Keyword("list")
	Dir     = data.Keyword("dir")
	File    = data.Keyword("file")

	OpenKey    = data.Keyword("open")
	ReadString = data.Keyword("read-string") // String
	ReadLines  = data.Keyword("read-lines")  // String
	ReadAll    = data.Keyword("read-all")    // Bytes
	ReadBlocks = data.Keyword("read-blocks") // Bytes
)

const defaultBlockSize = 4096

const (
	ErrExpectedDirectory   = "expected a directory, got a file: %s"
	ErrUnknownOpenMode     = "unknown open mode: %s"
	ErrExpectedFile        = "expected a file, got a directory: %s"
	ErrUnexpectedArguments = "unexpected additional arguments: %s"
)

var fileSystemType = types.MakeBasic("file-system")

func WrapFileSystem(fileSystem fs.FS) *FileSystem {
	return &FileSystem{
		fs:   fileSystem,
		List: bindList(fileSystem),
		Open: bindOpen(fileSystem),
	}
}

func (f *FileSystem) Type() types.Type {
	return fileSystemType
}

// Get returns the value associated with the specified key, if it exists.
func (f *FileSystem) Get(key data.Value) (data.Value, bool) {
	switch key {
	case ListKey:
		return f.List, true
	case OpenKey:
		return f.Open, true
	default:
		return nil, false
	}
}

func (f *FileSystem) Equal(other data.Value) bool {
	return f == other
}

func bindList(fileSystem fs.FS) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
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
	}, 1)
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

func bindOpen(fs fs.FS) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		path := args[0].(data.String)
		f, err := openFile(fs, path.String())
		if err != nil {
			panic(err)
		}
		return createReader(f, args[1:]...)
	}, 1, 3)
}

func createReader(f fs.File, args ...data.Value) data.Value {
	if len(args) == 0 {
		return NewReader(f, RuneInput)
	}
	switch args[0] {
	case ReadAll:
		return readAll(f)
	case ReadBlocks:
		size := getBlockSize(defaultBlockSize, args[1:]...)
		input, err := BlockInput(size)
		if err != nil {
			panic(err)
		}
		return NewReader(f, input)
	case ReadString:
		b := readAll(f)
		s := unsafe.String(&b[0], len(b))
		return data.String(s)
	case ReadLines:
		return NewReader(f, LineInput)
	default:
		panic(fmt.Errorf(ErrUnknownOpenMode, args[0]))
	}
}

func getBlockSize(def int, args ...data.Value) int {
	switch len(args) {
	case 0:
		return def
	case 1:
		return int(args[0].(data.Integer))
	default:
		panic(fmt.Errorf(ErrUnexpectedArguments, args[1:]))
	}
}

func openFile(fs fs.FS, path string) (fs.File, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if s.IsDir() {
		_ = f.Close()
		return nil, fmt.Errorf(ErrExpectedFile, path)
	}
	return f, nil
}

func readAll(f fs.File) data.Bytes {
	defer func() { _ = f.Close() }()
	size, err := getFileSize(f)
	if err != nil {
		panic(err)
	}
	buf := make(data.Bytes, size)
	l, err := f.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf[:l:l]
}

func getFileSize(f fs.File) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

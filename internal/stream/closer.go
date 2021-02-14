package stream

import "github.com/kode4food/ale/data"

// Closer is used to close a File
type Closer interface {
	Close()
}

const (
	// CloseKey is the key used to close a file
	CloseKey = data.Keyword("close")
)

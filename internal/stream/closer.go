package stream

import "github.com/kode4food/ale/data"

// Closer is used to close a File
type Closer interface {
	Close()
}

// CloseKey is the key used to close a file
const CloseKey = data.Keyword("close")

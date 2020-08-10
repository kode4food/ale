package stream

// Closer is used to close a File
type Closer interface {
	Close()
}

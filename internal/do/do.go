package do

import "sync"

// Action is a callback interface for eventually triggering some action
type Action func(func())

// Once creates a Do instance for performing an action only once
func Once() Action {
	var once sync.Once
	return once.Do
}

// Always returns a Do instance for always performing an action
func Always() Action {
	return func(f func()) {
		f()
	}
}

// Never returns a Do instance for never performing an action
func Never() Action {
	return func(func()) {
		// no-op
	}
}

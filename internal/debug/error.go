package debug

import (
	"errors"
	"fmt"
)

// ProgrammerError is raised when a condition arises that absolutely should not
// happen unless one of the compiler authors screwed up royally. Represented as
// a string to distinguish it from proper errors.
func ProgrammerError(msg string) string {
	return errors.New(msg).Error()
}

// ProgrammerErrorf is a formatted version of ProgrammerError.
func ProgrammerErrorf(format string, a ...any) string {
	return fmt.Errorf(format, a...).Error()
}

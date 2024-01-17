package debug

import "fmt"

// ProgrammerError is raised when a condition arises that absolutely should not
// happen unless one of the compiler authors screwed up royally. Represented as
// a string to distinguish it from proper errors.
func ProgrammerError(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

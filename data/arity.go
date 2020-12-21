package data

import "fmt"

// Error messages
const (
	errBadRangedArity = "got %d arguments, expected between %d and %d"
)

func checkRangedArity(min, max, count int) error {
	if count < min || count > max {
		return fmt.Errorf(errBadRangedArity, count, min, max)
	}
	return nil
}

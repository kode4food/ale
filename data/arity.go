package data

import "fmt"

const badRangedArity = "got %d arguments, expected between %d and %d"

func checkRangedArity(min, max, count int) error {
	if count < min || count > max {
		return fmt.Errorf(badRangedArity, count, min, max)
	}
	return nil
}

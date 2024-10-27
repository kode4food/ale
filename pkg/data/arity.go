package data

import (
	"fmt"

	"github.com/kode4food/ale/internal/debug"
)

const (
	// ErrFixedArity is raised when there are too few or many arguments
	// provided to a fixed ArityChecker
	ErrFixedArity = "expected %d arguments, got %d"

	// ErrMinimumArity is raised when there are too few arguments provided to a
	// minimum ArityChecker
	ErrMinimumArity = "expected at least %d arguments, got %d"

	// ErrRangedArity is raised when there are too few or many arguments
	// provided to a ranged ArityChecker
	ErrRangedArity = "expected between %d and %d arguments, got %d"

	// ErrTooManyArguments is raised when there are too many arguments provided
	// to a maximum ArityChecker
	ErrTooManyArguments = "too many arity check arguments"
)

// OrMore is the constant used when you want to tell MakeChecker to generate a
// minimum arity checker
const OrMore = -1

// MakeChecker produces an arity checker based on its parameters
func MakeChecker(arity ...int) ArityChecker {
	al := len(arity)
	switch {
	case al == 0:
		return CheckAnyArity
	case al > 2:
		panic(debug.ProgrammerError(ErrTooManyArguments))
	case al == 1 || arity[0] == arity[1]:
		return MakeFixedChecker(arity[0])
	case al == 2 && arity[1] == OrMore:
		return MakeMinimumChecker(arity[0])
	default:
		return MakeRangedChecker(arity[0], arity[1])
	}
}

// CheckAnyArity allows for any number of arguments
func CheckAnyArity(int) error {
	return nil
}

// MakeFixedChecker generates a fixed arity checker
func MakeFixedChecker(fixed int) ArityChecker {
	return func(count int) error {
		return CheckFixedArity(fixed, count)
	}
}

// CheckFixedArity allows for a fixed number of arguments
func CheckFixedArity(fixed, count int) error {
	if count != fixed {
		return fmt.Errorf(ErrFixedArity, fixed, count)
	}
	return nil
}

// MakeMinimumChecker generates a minimum arity checker
func MakeMinimumChecker(min int) ArityChecker {
	return func(count int) error {
		return CheckMinimumArity(min, count)
	}
}

// CheckMinimumArity allows for a minimum number of arguments
func CheckMinimumArity(min, count int) error {
	if count < min {
		return fmt.Errorf(ErrMinimumArity, min, count)
	}
	return nil
}

// MakeRangedChecker generates a ranged arity checker
func MakeRangedChecker(min, max int) ArityChecker {
	return func(count int) error {
		return CheckRangedArity(min, max, count)
	}
}

// CheckRangedArity allows for a ranged number of arguments
func CheckRangedArity(min, max, count int) error {
	if count < min || count > max {
		return fmt.Errorf(ErrRangedArity, min, max, count)
	}
	return nil
}

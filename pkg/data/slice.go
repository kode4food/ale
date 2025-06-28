package data

import "fmt"

func sliceRangedCall[T any](s []T, args Vector) ([]T, error) {
	switch len(args) {
	case 1:
		start := int(args[0].(Integer))
		if start < 0 || start > len(s) {
			return nil, fmt.Errorf(ErrInvalidStartIndex, start)
		}
		return s[start:], nil
	case 2:
		start := int(args[0].(Integer))
		end := int(args[1].(Integer))
		if start < 0 || end < start || end > len(s) {
			return nil, fmt.Errorf(ErrInvalidIndexes, start, end)
		}
		return s[start:end], nil
	default:
		return nil, fmt.Errorf(ErrRangedArity, 1, 2, len(args))
	}
}

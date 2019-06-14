package util

// IntMin returns the lesser of two integers
func IntMin(left, right int) int {
	if left < right {
		return left
	}
	return right
}

// IntMax returns the greater of two integers
func IntMax(left, right int) int {
	if left > right {
		return left
	}
	return right
}

package util

// Contain checks whether a element is in a slice (same type)
func Contain[T comparable](item T, slice []T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Reverse reverses any slice
// Source: https://stackoverflow.com/questions/28058278
func Reverse[T any](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

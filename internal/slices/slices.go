package slices

func MapIndex[S any, T any](slice []S, fn func(int, S) T) []T {
	result := make([]T, len(slice))
	for i, item := range slice {
		result[i] = fn(i, item)
	}
	return result
}

func CastSlice[T any, S any](slice []S, fn func(S) T) []T {
	result := make([]T, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

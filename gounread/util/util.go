package util

func Filter[T any](ts []T, fn func(T) bool) []T {
	result := make([]T, 0, len(ts))
	for _, t := range ts {
		if fn(t) {
			result = append(result, t)
		}
	}
	return result
}

func ForEach[T any](ts []T, fn func(T)) {
	for _, t := range ts {
		fn(t)
	}
}

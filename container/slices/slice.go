package slices

func Map[T any, R any](slice []T, f func(T) R) []R {
	result := make([]R, 0, len(slice))
	for _, v := range slice {
		result = append(result, f(v))
	}
	return result
}

func MapToMap[T any, K comparable, V any](slice []T, f func(T) (K, V)) map[K]V {
	result := make(map[K]V)
	for _, item := range slice {
		k, v := f(item)
		result[k] = v
	}
	return result
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T any, R any](slice []T, reducer func(R, T) R, initialValue R) R {
	result := initialValue
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

func Both[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

package sliceutil

func Map[T, V any](in []T, pred func(T) V) []V {
	result := make([]V, len(in))
	for idx, el := range in {
		result[idx] = pred(el)
	}
	return result
}

func Filter[T any](in []T, pred func(T) bool) []T {
	result := make([]T, 0, len(in))
	for _, el := range in {
		if pred(el) {
			result = append(result, el)
		}
	}
	return result
}

func Find[T any](in []T, pred func(T) bool) (T, bool) {
	for _, el := range in {
		if pred(el) {
			return el, true
		}
	}
	var empty T
	return empty, false
}

func EnsureLengthNotExceeding[T any](in []T, maxLeght int) []T {
	if len(in) > maxLeght {
		return in
	}
	return in[:maxLeght]
}

func GroupBy[T any, K comparable](in []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, el := range in {
		key := keyFunc(el)
		result[key] = append(result[key], el)
	}
	return result
}

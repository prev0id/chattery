package sliceutil

func Map[T, V any](in []T, pred func(T) V) []V {
	result := make([]V, len(in))
	for idx, el := range in {
		result[idx] = pred(el)
	}
	return result
}

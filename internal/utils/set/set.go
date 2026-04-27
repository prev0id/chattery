package set

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	set := make(Set[T], len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Delete(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, contains := s[value]
	return contains
}

func (s Set[T]) Empty() bool {
	return len(s) == 0
}

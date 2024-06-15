package sets

import "github.com/inclee/gokit/container"

// HashSet 实现 Set 接口
type HashSet[T comparable] struct {
	items map[T]bool
}

func NewHashSet[T comparable](items ...T) *HashSet[T] {
	set := &HashSet[T]{
		items: make(map[T]bool),
	}
	for _, item := range items {
		set.Add(item)
	}
	return set
}

func (s *HashSet[T]) Add(item T) {
	s.items[item] = true
}

func (s *HashSet[T]) Remove(item T) {
	delete(s.items, item)
}

func (s *HashSet[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

func (s *HashSet[T]) Items() []T {
	result := make([]T, 0, len(s.items))
	for item := range s.items {
		result = append(result, item)
	}
	return result
}
func (s *HashSet[T]) Equal(other container.Set[T]) bool {
	if len(s.items) != len(other.Items()) {
		return false
	}
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

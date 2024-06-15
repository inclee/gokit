package container

// Set 接口定义
type Set[T comparable] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Items() []T
	Equal(other Set[T]) bool
}

package maps

import (
	"github.com/inclee/gokit/container"
	"github.com/inclee/gokit/container/sets"
)

func Keys[K comparable, V any](m map[K]V) container.Set[K] {
	keys := sets.NewHashSet[K]()
	for k := range m {
		keys.Add(k)
	}
	return keys
}

func Equal[K comparable, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}
	for key, va := range a {
		vb, exists := b[key]
		if !exists || va != vb {
			return false
		}
	}
	return true
}

func Both[K comparable, V any](a, b map[K]V, fn func(K, V, V) bool) bool {
	if len(a) < len(b) {
		return false
	}
	for key, va := range a {
		vb, exists := b[key]
		if !exists || !fn(key, va, vb) {
			return false
		}
	}
	return true
}

func Any[K comparable, V any](a, b map[K]V, fn func(K, V, V) bool) bool {
	for key, va := range a {
		vb, exists := b[key]
		if exists || fn(key, va, vb) {
			return true
		}
	}
	return false
}

// Merge Map b into Map a, overriding values of keys that exist when 'override' is true.
func Merge[K comparable, V any](a, b map[K]V, override bool) {
	for bk, bv := range b {
		if _, ok := a[bk]; ok && override {
			a[bk] = bv
			continue
		}
		a[bk] = bv
	}
}

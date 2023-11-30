package set

import "fediverse/iterable"

// Set is for representing a set of objects, irrespective insertion order.
// additionally, duplicate insertion of the same key into a Set will result in
// subsequent `Add` invocations to effectively be a no-op
type Set[K comparable] map[K]bool

var _ iterable.Iterable[string] = Set[string]{}

func New[K comparable](items ...K) Set[K] {
	result := Set[K]{}

	for _, item := range items {
		result.Add(item)
	}

	return result
}

// FromSlice creates a new Set from an existing slice
func FromSlice[K comparable](s []K) Set[K] {
	newSet := Set[K]{}
	for _, v := range s {
		newSet.Add(v)
	}

	return newSet
}

// Add adds a key to the set
func (s Set[K]) Add(value K) {
	s[value] = true
}

// Has checks for existence of a key in the set
func (s Set[K]) Has(value K) bool {
	return s[value]
}

// Iterate creates a channel purely for iteration purposes
func (s Set[K]) Iterate() <-chan K {
	c := make(chan K)
	go func() {
		for k := range s {
			c <- k
		}
		close(c)
	}()

	return c
}

// Equal checks two set equality.
//
// What's nice about sets is that insertion order does not matter; only the
// cardinality of the set, and also wither all existing keys match each other's
// keys will be checked against
func (s Set[K]) Equals(s1 Set[K]) bool {
	if len(s) != len(s1) {
		return false
	}

	for k := range s {
		if !s1.Has(k) {
			return false
		}
	}

	return true
}

// IsSubsetTo checks if the current set is a subset of another set
func (s Set[K]) IsSubsetTo(s1 Set[K]) bool {
	if len(s) > len(s1) {
		return false
	}

	for k := range s {
		if !s1.Has(k) {
			return false
		}
	}

	return true
}

// Union performs a union of two sets. You can think of this method as if it
// were to concatenate two sets together
func (s Set[K]) Union(s1 Set[K]) Set[K] {
	result := Set[K]{}
	for k, ok := range s {
		if ok {
			result.Add(k)
		}
	}
	for k, ok := range s1 {
		if ok {
			result.Add(k)
		}
	}
	return result
}

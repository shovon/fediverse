package memoizer

type Cache[K comparable, V any] struct {
	cache map[K]V
}

func (c Cache[K, V]) getCache() map[K]V {
	if c.cache == nil {
		c.cache = make(map[K]V)
	}

	return c.cache
}

func (c Cache[K, V]) Memoize(input K, d func(input K) V) V {
	cache := c.getCache()

	if _, ok := cache[input]; !ok {
		cache[input] = d(input)
	}

	return cache[input]
}

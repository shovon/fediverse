package slices

func Reduce[V, T any](s []V, initial T, a func(result T, next V, index int) T) T {
	accum := initial
	for i, value := range s {
		accum = a(accum, value, i)
	}
	return accum
}

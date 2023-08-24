package functional

func Compose[T any, V any, Y any](g func(V) Y, f func(T) V) func(T) Y {
	return func(t T) Y {
		return g(f(t))
	}
}

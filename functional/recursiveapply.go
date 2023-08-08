package functional

func RecursiveApply[T any](f [](func(T) T)) func(T) T {
	if len(f) == 0 {
		return func(d T) T {
			return d
		}
	}
	if len(f) == 1 {
		return f[0]
	}

	return func(d T) T {
		return f[0](RecursiveApply(f[1:])(d))
	}
}

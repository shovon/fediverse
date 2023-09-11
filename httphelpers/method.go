package httphelpers

func Method(method string) Processor {
	return Condition(func(r ReadOnlyRequest) bool {
		return r.Method == method
	})
}

package httphelpers

func Method(method string) Processor {
	return Condition(func(r BarebonesRequest) bool {
		return r.Method == method
	})
}

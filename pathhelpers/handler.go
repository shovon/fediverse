package pathhelpers

// TODO: unit test this

type Handlers map[string](func(params map[string]string) (bool, map[string]string))

func (h Handlers) Handle(path string) (bool, map[string]string) {
	for key, fn := range h {
		match, params := Match(key, path)
		if match {
			return fn(params)
		}
	}

	return false, nil
}

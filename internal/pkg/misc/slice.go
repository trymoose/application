package misc

func MapSlice[I, O any](i []I, fn func(I) O) []O {
	o := make([]O, len(i))
	for idx, elem := range i {
		o[idx] = fn(elem)
	}
	return o
}

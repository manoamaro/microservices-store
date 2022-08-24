package collections

func MapTo[I interface{}, O interface{}](i []I, f func(I) O) []O {
	var output []O
	for _, el := range i {
		output = append(output, f(el))
	}
	return output
}

func Contains[T comparable](a []T, v T) bool {
	for _, t := range a {
		if t == v {
			return true
		}
	}
	return false
}

func ContainsAny[T comparable](a []T, b []T) bool {
	for _, va := range a {
		for _, vb := range b {
			if va == vb {
				return true
			}
		}
	}
	return false
}

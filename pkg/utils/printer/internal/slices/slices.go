package slices

func Map[T any, V any](v []T, fn func(T) V) []V {
	var out []V
	for _, v := range v {
		out = append(out, fn(v))
	}
	return out
}

func Filter[T any](v []T, fn func(T) bool) []T {
	var out []T
	for _, v := range v {
		if fn(v) {
			out = append(out, v)
		}
	}
	return out
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

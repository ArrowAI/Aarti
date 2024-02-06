package slices

func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func MapErr[T, U any](s []T, f func(T) (U, error)) ([]U, error) {
	r := make([]U, len(s))
	for i, v := range s {
		var err error
		r[i], err = f(v)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func Distinct[T comparable](s []T) []T {
	m := make(map[T]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	r := make([]T, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

package storage

func defaults[T comparable](v T, d T) T {
	var z T
	if v == z {
		return d
	}
	return v
}

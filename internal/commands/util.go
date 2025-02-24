package commands

func Optional[T comparable](v T) *T {
	var zero T

	if v != zero {
		return &v
	}

	return nil
}

func Ptr[T any](t T) *T {
	return &t
}

func Val[T any](t *T) T {
	if t == nil {
		var zero T

		return zero
	}

	return *t
}

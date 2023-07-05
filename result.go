package result

type Result[T any] struct {
	value T
	err   error
}

func (r Result[T]) Value() T {
	return r.value
}

func (r Result[T]) Error() error {
	return r.err
}

func New[T any](value T, rError error) Result[T] {
	if rError != nil {
		return err[T](rError)
	}

	return ok[T](value)
}

func (r Result[T]) Unwrap() (T, error) {
	return r.Value(), r.Error()
}

func ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

package result

type mFunc[T any, U any] func(T) U

// Result is a container type that holds a value of type T or an error.
// It cannot simultaneously hold a non-zero value and a non-nil error.
//
// Nevertheless, when working with Result, if it holds both an error and a value,
// it will always return a zero value of type T.
// In such cases, for example, calling Value() will return a zero value of type T,
// and calling Match() will always return result.IsErr().
type Result[T any, E any] struct {
	value  *T
	err    *E
	strict bool
}

// C is an optional config struct used to create a new Result.
type C struct {
	// Strict allows you to deny access to a Result value without using
	// the Error() method first.
	//
	// Calling Value() or Unwrap() on Result will panic whether there is
	// and error or not.
	Strict bool
}

// Ok returns a new instance of Result[T] that holds a value, where T is
// the inherited type of value.
//
// Calling Match() on such a Result will always match result.IsOk()
func Ok[T any](value T, config ...C) Result[T, T] {
	return ok[T, T](value, config...)
}

// Err returns a new instance of Result[T] that holds an error, where T is
// the inherited type of zero value.
//
// Calling Match() on such a Result will always match result.IsErr()
func Err[E any](rErr E, config ...C) Result[E, E] {
	return err[E, E](rErr, config...)
}

// newResult returns an instance of the Result[T], where T is
// the inherited type of value.
//
// If a non-nil rError was passed as an argument, Error() will return
// the provided error, and Value() will return a zero value of type T.
//
// Otherwise, Value() will return the valid value passed earlier,
// and Error() will return nil instead of an error.
//
// The discriminated union nature of T | error is simulated using the Match() method,
// which allows matching the current Result[T] with result.IsOk() in case of no error,
// as well as with result.IsErr() in case of an error.
func newResult[T, E any](value T, rError E, config ...C) Result[T, E] {
	var isStrict bool

	if len(config) > 0 {
		isStrict = config[0].Strict
	}

	return Result[T, E]{value: &value, err: &rError, strict: isStrict}
}

// Match simulates pattern matching for handling scenarios of having a valid result
// or having an error.
// It is used in the following format:
//
//	v := 1
//	r := result.New(1, nil)
//
//	switch r.Match() {
//	case result.IsOk():
//		fmt.Println("there is a value:")
//		fmt.Println(x.Value())
//	case result.IsErr():
//		fmt.Println("there is an error:")
//		fmt.Println(x.Error())
//	}
func Match[T, E any, U any](r Result[T, E], okF func(T) U, errF func(E) U) U {
	if r.err != nil {
		return errF(r.innerError())
	}

	return okF(r.innerValue())
}

func (r Result[T, E]) innerValue() T {
	return *r.value
}

func (r Result[T, E]) innerError() E {
	// if r.err != nil {
	// 	return *r.err
	// }

	// return *new(E)
	return *r.err
}

// And returns a passed newR if the Result has no error and newR has no error.
//
// Otherwise returns a Result with a contained error, if present,
// or a new Result with error from newR.
func (r Result[T, E]) And(newR Result[T, E]) Result[T, E] {
	if r.err != nil {
		return err[T, E](r.innerError())
	}

	if newR.err != nil {
		return err[T, E](newR.innerError())
	}

	return ok[T, E](newR.innerValue())
}

// AndThen returns the result of f function with Result value as an argument if
// Result has no error.
//
// Otherwise returns a Result with a contained error.
func (r Result[T, E]) AndThen(f func(T) Result[T, E]) Result[T, E] {
	if r.err != nil {
		return err[T, E](r.innerError())
	}

	return f(r.innerValue())
}

// IsErr returns true if Result has an error.
func (r Result[T, E]) IsErr() bool {
	return r.err != nil
}

// IsErrAnd returns true if rErr matches the contained Result error.
func (r Result[T, E]) IsErrAnd(rErr E) bool {
	if r.err == nil {
		return false
	}

	return r.err == &rErr
}

// IsOk returns true if Result has no error.
func (r Result[T, E]) IsOk() bool {
	return r.err == nil
}

// IsOkAnd returns true if Result has no error and the contained value
// matches a predicate of f.
func (r Result[T, E]) IsOkAnd(f func(T) bool) bool {
	if r.err != nil {
		return false
	}

	return f(r.innerValue())
}

// Map returns a new Result by applying an f function to a Result value,
// leaving Result error untouched.
func Map[T, E, U any](r Result[T, E], okF mFunc[T, U]) Result[U, E] {
	computed := okF(r.innerValue())

	return Result[U, E]{value: &computed, err: r.err}
}

// MapErr returns result of f function with Result error as an argument
// if Result has an error.
//
// Otherwise returns self.
func (r Result[T, E]) MapErr(errF func(E) E) Result[T, E] {
	if r.err == nil {
		return r
	}

	computed := errF(r.innerError())

	return Result[T, E]{value: r.value, err: &computed}
}

// MapOr returns the provided rDefault of type T if Result has an error, otherwise
// returns a result of f function with Result value as an argument.
func (r Result[T, E]) MapOr(rDefault T, f func(T) T) T {
	if r.err != nil {
		return rDefault
	}

	return f(r.innerValue())
}

// MapOrElse calls an errF function with Result error as an argument if Result
// has an error, otherwise calls an okF function with a Result value as an argument.
func (r Result[T, E]) MapOrElse(errF func(E) T, okF func(T) T) T {
	if r.err != nil {
		return errF(r.innerError())
	}

	return okF(r.innerValue())
}

// Unwrap allows obtaining the nested value and error inside the Result as
// a (T, error) return signature.
//
// It is recommended to use the Match() method for proper pattern matching.
//
// Will panic if Result was built with C{strict: true}
func (r Result[T, E]) Unwrap() (T, E) {
	if r.strict {
		panic("cannot unwrap value and error of strict result")
	}
	return r.innerValue(), r.innerError()
}

func (r Result[T, E]) UnwrapOr(rDefault T) T {
	if r.err != nil {
		return rDefault
	}

	return r.innerValue()
}

func (r Result[T, E]) UnwrapOrElse(f func(E) T) T {
	if r.err == nil {
		return r.innerValue()
	}

	return f(r.innerError())
}

// IsStrict returns whether calling Unwrap() will panic or not.
func (r Result[T, E]) IsStrict() bool {
	return r.strict
}

func ok[T, E any](value T, config ...C) Result[T, E] {
	var isStrict bool

	if len(config) > 0 {
		isStrict = config[0].Strict
	}

	return Result[T, E]{value: &value, strict: isStrict}
}

func err[T, E any](err E, config ...C) Result[T, E] {
	var isStrict bool

	if len(config) > 0 {
		isStrict = config[0].Strict
	}

	return Result[T, E]{err: &err, strict: isStrict}
}

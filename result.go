package result

import (
	"fmt"
)

type mFunc[T any, U any] func(T) U

// Result is a container type that holds a value of type T or an error.
// It cannot simultaneously hold a non-zero value and a non-nil error.
//
// Nevertheless, when working with Result, if it holds both an error and a value,
// it will always return a zero value of type T.
// In such cases, for example, calling Value() will return a zero value of type T,
// and calling Match() will always return result.IsErr().
type Result[T any, U any] struct {
	value *T
	err   error
}

// Ok returns a new instance of Result[T, any] that holds a value, where T is
// the inherited type of value.
//
// Calling Match() on such a Result will always call an okF function.
func Ok[T any](value T) Result[T, any] {
	return ok[T](value)
}

// Err returns a new instance of Result[T, any] that holds an error, where T is
// the inherited type of error.
//
// Calling Match() on such a Result will always call an errF function.
func Err[T any](rErr error) Result[T, any] {
	return err[T](rErr)
}

// newResult returns an instance of the Result[T, any], where T is
// the inherited type of value.
func newResult[T any](value T, rError error) Result[T, any] {
	return Result[T, any]{value: &value, err: rError}
}

// Match simulates pattern matching for handling scenarios of having a valid result
// or having an error.
// It is used in the following format:
//
//	userId := result.Ok(1337)
//
//	res := result.Match(userId,
//		func(v int) int { return v * 2 },
//		func(err int) int { return err * 0 },
//	)
//
//	fmt.Println(res)
func Match[T any, U any](r Result[T, any], okF mFunc[T, U], errF func(error) U) U {
	if r.isErr() {
		return errF(r.innerError())
	}

	return okF(r.innerValue())
}

func (r Result[T, any]) innerValue() T {
	if r.value == nil {
		return *new(T)
	}

	return *r.value
}

func (r Result[T, any]) innerError() error {
	if r.err == nil {
		return nil
	}

	return r.err.(error)
}

func (r Result[T, any]) isErr() bool {
	return r.err != nil
}

func (r Result[T, any]) isOk() bool {
	return r.err == nil
}

// And returns a passed newR if the Result has no error and newR has no error.
//
// Otherwise returns a Result with a contained error, if present,
// or a new Result with error from newR.
func (r Result[T, any]) And(newR Result[T, any]) Result[T, any] {
	if r.isErr() {
		return r
	}

	if newR.isErr() {
		return newR
	}

	return r
}

// AndThen returns the result of f function with Result value as an argument if
// Result has no error.
//
// Otherwise returns a Result with a contained error.
func (r Result[T, any]) AndThen(f func(T) Result[T, any]) Result[T, any] {
	if r.isErr() {
		return r
	}

	return f(r.innerValue())
}

// IsErr returns true if Result has an error.
func (r Result[T, any]) IsErr() bool {
	return r.isErr()
}

// IsErrAnd returns true if rErr matches the contained Result error.
func (r Result[T, any]) IsErrAnd(f func(error) bool) bool {
	if r.isOk() {
		return false
	}

	return f(r.innerError())
}

// IsOk returns true if Result has no error.
func (r Result[T, any]) IsOk() bool {
	return r.isOk()
}

// IsOkAnd returns true if Result has no error and the contained value
// matches a predicate of f.
func (r Result[T, any]) IsOkAnd(f func(T) bool) bool {
	if r.isErr() {
		return false
	}

	return f(r.innerValue())
}

// Map returns a new Result by applying an f function to a Result value,
// leaving Result error untouched.
func Map[T, U any](r Result[T, any], okF mFunc[T, U]) Result[U, any] {
	computedValue := okF(r.innerValue())

	return newResult(computedValue, r.innerError())
}

func Expand[U, T any](r Result[T, any]) Result[T, U] {
	value := r.innerValue()

	return Result[T, U]{value: &value, err: r.innerError()}
}

func (r Result[T, U]) Map(okF mFunc[T, U]) Result[U, any] {
	computedValue := okF(r.innerValue())

	return Result[U, any]{value: &computedValue, err: r.innerError()}
}

// MapErr returns result of errF function with Result error as an argument
// if Result has an error.
//
// Otherwise returns self.
func (r Result[T, any]) MapErr(errF func(error) error) Result[T, any] {
	if r.isOk() {
		return r
	}

	computedErr := errF(r.innerError())
	value := r.innerValue()

	return Result[T, any]{value: &value, err: computedErr}
}

// MapOr returns the provided rDefault of type T if Result has an error, otherwise
// returns a result of f function with Result value as an argument.
func (r Result[T, any]) MapOr(rDefault T, f func(T) T) T {
	if r.isErr() {
		return rDefault
	}

	return f(r.innerValue())
}

// MapOrElse calls an errF function with Result error as an argument if Result
// has an error, otherwise calls an okF function with a Result value as an argument.
func (r Result[T, any]) MapOrElse(errF func(error) T, okF func(T) T) T {
	if r.isErr() {
		return errF(r.innerError())
	}

	return okF(r.innerValue())
}

func (r Result[T, any]) Expect(msg string) Result[T, any] {
	if r.isOk() {
		return r
	}

	wrappedErr := fmt.Errorf("%v: %w", msg, r.err)
	value := r.innerValue()

	return Result[T, any]{value: &value, err: wrappedErr}
}

// Unwrap allows obtaining the nested value and error inside the Result as
// a (T, error) return signature.
//
// It is recommended to use the Match() method for proper pattern matching.
//
// Will panic if Result was built with C{strict: true}
func (r Result[T, any]) Unwrap() (T, error) {
	return r.innerValue(), r.innerError()
}

func (r Result[T, any]) UnwrapOr(rDefault T) T {
	if r.isOk() {
		return r.innerValue()
	}

	return rDefault
}

func (r Result[T, any]) UnwrapOrElse(f func(error) T) T {
	if r.isOk() {
		return r.innerValue()
	}

	return f(r.innerError())
}

func ok[T any](value T) Result[T, any] {
	return Result[T, any]{value: &value}
}

func err[T any](err error) Result[T, any] {
	return Result[T, any]{err: err}
}

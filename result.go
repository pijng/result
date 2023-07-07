package result

import (
	"errors"
	"fmt"
	"reflect"
)

type isOk bool
type isErr bool

type mFunc[T any, U any] func(T) U

// Result is a container type that holds a value of type T or an error.
// It cannot simultaneously hold a non-nil value and a non-nil error.
//
// Nevertheless, when working with Result, if it holds both an error and a value,
// it will always return a nil value of type T.
// In such cases, for example, calling Value() will return a nil value of type T,
// and calling Match() will always return result.Err().
type Result[T any] struct {
	value   *T
	err     error
	variant interface{}
}

// New returns a new instance of the Result[T] type, where T is
// the inherited type of value.
//
// If a non-nil rError was passed as an argument, Error() will return
// the provided error, and Value() will return a nil value value of type T.
//
// Otherwise, Value() will return the valid value passed earlier,
// and Error() will return nil instead of an error.
//
// The discriminated union nature of T | error is simulated using the Match() method,
// which allows matching the current Result[T] with result.Ok() in case of no error,
// as well as with result.Err() in case of an error.
func New[T any](value T, rError error) Result[T] {
	if rError != nil {
		return err[T](rError)
	}

	return ok[T](value)
}

// Match simulates pattern matching for handling scenarios of having a valid result
// or having an error.
// It is used in the following format:
//
//	v := 1
//	r := result.New(1, nil)
//	switch r.Match() {
//	case result.Ok():
//		fmt.Println("there is a value:")
//		fmt.Println(x.Value())
//	case result.Err():
//		fmt.Println("there is an error:")
//		fmt.Println(x.Error())
//	}
func (r Result[T]) Match() reflect.Type {
	return reflect.TypeOf(r.variant)
}

// Value returns the underlying value of type T.
func (r Result[T]) Value() T {
	if r.value == nil {
		return *new(T)
	}

	return *r.value
}

// Error returns the underlying error.
func (r Result[T]) Error() error {
	return r.err
}

// And returns a passed newR if the Result has no error and newR has no error.
//
// Otherwise returns a Result with a contained error, if present,
// or a new Result with error from newR.
func (r Result[T]) And(newR Result[T]) Result[T] {
	if r.err != nil {
		return err[T](r.err)
	}

	if newR.err != nil {
		return err[T](newR.err)
	}

	return ok[T](newR.Value())
}

// AndThen returns the result of f function with Result value as an argument if
// Result has no error.
//
// Otherwise returns a Result with a contained error.
func (r Result[T]) AndThen(f func(T) Result[T]) Result[T] {
	if r.err != nil {
		return err[T](r.err)
	}

	return f(r.Value())
}

// Expect returns a new Result by wrapping a contained error with rErr if error is present.
//
// Otherwise returns self.
func (r Result[T]) Expect(rErr error) Result[T] {
	if r.err != nil {
		return err[T](fmt.Errorf("%w: %w", rErr, r.err))
	}

	return r
}

// IsErr returns true if Result has an error.
func (r Result[T]) IsErr() bool {
	return r.Error() != nil
}

// IsErrAnd returns true if rErr matches the contained Result error.
func (r Result[T]) IsErrAnd(rErr error) bool {
	return errors.Is(r.Error(), rErr)
}

// IsOk returns true if Result has no error.
func (r Result[T]) IsOk() bool {
	return r.Error() == nil
}

// IsOkAnd returns true if Result has no error and the contained value
// matches a predicate of f.
func (r Result[T]) IsOkAnd(f func(T) bool) bool {
	if r.Error() != nil {
		return false
	}

	return f(r.Value())
}

// Map returns a new Result by applying an f function to a Result value,
// leaving Result error untouched.
func Map[T, U any](r Result[T], f mFunc[T, U]) Result[U] {
	return New[U](f(r.Value()), r.Error())
}

// MapErr returns result of f function with Result error as an argument
// if Result has an error.
//
// Otherwise returns self.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.Error() == nil {
		return r
	}

	return New(r.Value(), f(r.Error()))
}

// MapOr returns the provided v of type T if Result has an error, otherwise
// returns a result of f function with Result value as an argument.
func (r Result[T]) MapOr(v T, f func(T) T) T {
	if r.Error() != nil {
		return v
	}

	return f(r.Value())
}

// MapOrElse calls an eF function with Result error as an argument if Result
// has an error, otherwise calls an oF function with a Result value as an argument.
func (r Result[T]) MapOrElse(eF func(error) T, oF func(T) T) T {
	if r.Error() != nil {
		return eF(r.Error())
	}

	return oF(r.Value())
}

// Unwrap allows obtaining the nested value and error inside the Result as
// a (T, error) return signature.
//
// It is recommended to use the Match() method for proper pattern matching.
func (r Result[T]) Unwrap() (T, error) {
	return r.Value(), r.Error()
}

// Ok returns the type to match with Match() for valid cases.
func Ok() reflect.Type {
	return reflect.TypeOf(new(isOk))
}

// Err returns the type to match with Match() for invalid cases.
func Err() reflect.Type {
	return reflect.TypeOf(new(isErr))
}

func ok[T any](value T) Result[T] {
	return Result[T]{value: &value, variant: new(isOk)}
}

func err[T any](err error) Result[T] {
	return Result[T]{err: err, variant: new(isErr)}
}

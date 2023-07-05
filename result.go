package result

import (
	"fmt"
	"reflect"
)

type isOk bool
type isErr bool

type Result[T any] struct {
	value   T
	err     error
	variant interface{}
}

// New returns a new instance of the Result[T] type, where T is
// the inherited type of value.
//
// If a non-nil rError was passed as an argument, Result.Error() will return
// the provided error, and Result.Value() will contain a nil pointer.
//
// Otherwise, Result.Value() will return the valid value passed earlier,
// and Result.Error() will return nil instead of an error.
//
// The discriminated union nature of T | error is simulated using the Result.Match() method,
// which allows matching the current Result[T] with Result.Ok() in case of no error,
// as well as with Result.Err() in case of an error.
func New[T any](value T, rError error) Result[T] {
	if rError != nil {
		return err[T](rError)
	}

	return ok[T](value)
}

// Match simulates pattern matching for handling scenarios of having a valid result or having an error.
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

// Value returns the underlying value of the Result with type T.
func (r Result[T]) Value() T {
	return r.value
}

// Error returns the underlying error of the Result.
func (r Result[T]) Error() error {
	return r.err
}

// Ok returns the type to match with Result.Match() for valid cases.
func Ok() reflect.Type {
	return reflect.TypeOf(new(isOk))
}

// Err returns the type to match with Result.Match() for invalid cases.
func Err() reflect.Type {
	return reflect.TypeOf(new(isErr))
}

// Unwrap allows obtaining the nested value and error inside the Result as a (T, error) return signature.
//
// It is recommended to use the Result.Match() method for proper pattern matching.
//
// If an error is present, attempting to work with T may panic.
// If T is present, attempting to work with the error may panic.
func (r Result[T]) Unwrap() (T, error) {
	return r.Value(), r.Error()
}

func ok[T any](value T) Result[T] {
	return Result[T]{value: value, variant: new(isOk)}
}

func err[T any](err error) Result[T] {
	return Result[T]{err: err, variant: new(isErr)}
}

func test() {
	v := 1
	x := New(v, nil)

	switch x.Match() {
	case Ok():
		fmt.Println("there is value:")
		fmt.Println(x.Value())
	case Err():
		fmt.Println("there is error:")
		fmt.Println(x.Error())
	}
}

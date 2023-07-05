# result

`result` is a package that naively implements `Result<T, E>` from Rust in Go.

It is used for returning and propagating values and errors.

## Installation

Install with `go mod`

```bash
go get github.com/pijng/result
```

## Usage/Examples

#### Creating a new `Return` object and returning it as the function's result:

```go
type User struct {
  name string
  age  int
}

func getUser(id int) result.Result[User] {
  if id == 0 {
    return result.New(User{}, fmt.Errorf("user with id %d not found", id))
  }

  user := User{name: "pijng", age: 26}

  return result.New(user, nil)
}

func main() {
  r1 := getUser(1)
  r2 := getUser(0)

  fmt.Println(r1.Value()) // outputs: {pijng 26}
  fmt.Println(r1.Error()) // outputs: <nil>
  fmt.Println(r2.Value()) // outputs: { 0}
  fmt.Println(r2.Error()) // outputs: user with id 0 not found
}
```

#### Usage with pattern matching:

```go
type User struct {
  name string
  age  int
}

func getUser(id int) result.Result[User] {
  if id == 0 {
    return result.New(User{}, fmt.Errorf("user with id %d not found", id))
  }

  user := User{name: "pijng", age: 26}

  return result.New(user, nil)
}

func main() {
    r1 := getUser(1)

    switch r1.Match() {
    case result.Ok(): // This block will always be called.
        fmt.Println(r1.Value())
    case result.Err(): // This block will never be called.
        fmt.Println(r1.Error())
    }

    r2 := getUser(0)

    switch r2.Match() {
    case result.Ok(): // This block will never be called.
        fmt.Println(r2.Value())
    case result.Err(): // This block will always be called.
        fmt.Println(r2.Error())
    }
}
```

#### Unwrap value and error from `Result[T]`:

```go
func getUser(id int) result.Result[User] {
  if id == 0 {
    return result.New(User{}, fmt.Errorf("user with id %d not found", id))
  }

  user := User{name: "pijng", age: 26}

  return result.New(user, nil)
}

func main() {
    r1 := getUser(1)

    user, err := r1.Unwrap()

    fmt.Println(user) // outputs: {pijng 26}
    fmt.Println(err) // outputs: <nil>

    r2 := getUser(0)

    user, err := r2.Unwrap()

    fmt.Println(user) // outputs: { 0}
    fmt.Println(err) // outputs: user with id 0 not found
}

```

## FAQ

#### Why do that in the first place?

Honestly, it doesn't make much sense.

In the end, working with `Result[T]` boils down to either simulating pattern matching, which doesn't differ much from a regular `if err != nil`, or propagating data up the layers, which is also similar to the standard `return (T, error)` signature.

This package was created out of curiosity. While it serves its function of mimicking `Result<T, E>`, it doesn't bring any significant advantages in itself.

#### Can this package be used in production?

It's better not to. On one hand, everything works correctly, but on the other hand, you will introduce unnecessary complexity to the project.

But in general, you can :)

## License

[MIT](https://choosealicense.com/licenses/mit/)

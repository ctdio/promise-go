# promise-go

Simple, generic promises/futures with go.

### Installation

```bash
go get -u github.com/charlieduong94/promise-go
```

### Usage

Creating a promise is quite easy, all you need to do is pass in a function that
you want to run asynchronously into the `New` function, which returns a channel.

**Note:** The function needs to have no arguments and return a generic interface and an error.

```go
myFunc := func () (interface{}, error) {
  return true, nil
}
```

After creating that function, pass it into the package's `New` function. This will create a channel
that will produce a pointer to a `Result` object when the async function produces a value.

Here's what the `Result` struct looks like:

```go
type Result struct {
  Value interface{}
  Error error
}
```

Below is a example of how a promise is used.

```go
package main

import (
  "fmt"
  "github.com/charlieduong94/promise-go" // exported package name is "promise"
)

type myStruct struct {
  Message string
}

func main () {
  // pass in a function that returns a generic interface and an error
  myPromise := promise.New(func () (interface{}, error) {
    // do something you want handled async here
    time.Sleep(5 * time.Second)

    val := myStruct{"this is a return value"}

    // return any results
    return val, nil
  })

  // perform any other work here...

  // when you are ready, grab the result from your promise
  result := <-myPromise
  if result.Error != nil {
    // handle errors
  }

  // the result of your function can be accessed via the "Value" attribute
  fmt.Println(result.Value)

  // of course, you can cast the value back to whatever return type you need
  myVal, err := result.Value.(myStruct)
  if err != nil {
    // handle error
  }

  fmt.Println(myVal.Message)
}
```

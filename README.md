# promise-go

Simple, generic promises/futures with go.

### Installation

```bash
go get -u github.com/charlieduong94/promise-go
```

### Usage

Creating a promise is quite easy. All you need to do is pass in a function that
you want to run asynchronously into the `New` function. The function needs to return
some sort of value and an error.

```go
package main

import (
  "fmt"
  "github.com/charlieduong94/promises-go" // exported package name is "promise"
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
  if err == nil {
    // handle error
  }

  fmt.Println(myVal.Message)
}
```

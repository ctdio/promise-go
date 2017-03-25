# promise-go

Simple, generic promises/futures with go.

### Installation

```bash
go get -u github.com/charlieduong94/promise-go
```

### Usage

#### Simple Promises

Creating a promise is quite easy, first create a function that you want to run
asynchronously.

**Note:** The function needs to have no arguments and return a generic interface and an error.

```go
myFunc := func () (interface{}, error) {
  return true, nil
}
```

After creating that function, pass it into the package's `Create` function. This will produce
a pointer to a `Promise` object. This object contains a method `GetResult()`, which will
block until the promise is resolved.

Here's what the `Result` struct looks like:

```go
type Result struct {
  Value interface{}
  Error error
}
```

Below is a example of how a promise can be used.

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
  myPromise := promise.Create(func () (interface{}, error) {
    // do something you want handled async here
    time.Sleep(5 * time.Second)

    val := myStruct{"this is a return value"}

    // return any results
    return val, nil
  })

  // perform any other work here...

  // when you are ready, grab the result from your promise
  result := myPromise.GetResult()
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

#### Multiple concurrent promises

Sometimes, you need to kick off multiple concurrent functions off at the same time. With the `CreateAll`
function, you can do that easily. Much like the `Create` function, this returns a promise.

**Note:** When dealing with combined promises, any error that is returned causes the promise
to immediately return the result with the first error that occurs.


```go
multiplePromises := promise.CreateAll(
  func () (interface{}, error) {
    time.Sleep(1 * time.Second)
    return 26, nil
  },
  func () (interface{}, error) {
    time.Sleep(2 * time.Second)
    return 58, nil
  },
)

combinedResult := multiplePromises.GetResult()
if combinedResult.Error != nil {
  // handle errors
}

fmt.Println(combinedResult.Value) // prints "[26 58]"

// since "Value" is returned as an interface{},
// assert the value back to that of slice
values := combinedResult.Value.([]interface{})

fmt.Println(values[0]) // 26
fmt.Println(values[1]) // 58
```


There may also be times when you may want to start promises at different times, but then await for all
of them to be resolved at a later phase. This can be done with the `All` function.

```go
promiseA := promise.Create(func () (interface{}, error) {
  return "A", nil
})

// do some work additional work

promiseB := promise.Create(func () (interface{}, error) {
  return "B", nil
})

combinedPromise := promise.All(promiseA, promiseB)
combinedResult := combinedPromise.GetResult()

// The result can be handled like in previous example
```

If you want to get fancy, you can use `All` to take multiple Promises from `Create` and `CreateAll`
to await for them all at once.

```go
promiseA := promise.Create(func () (interface{}, error) {
  return "A", nil
})

// do some work additional work

promiseB := promise.CreateAll(
  func () (interface{}, error) {
    return "B", nil
  },
  func () (interface{}, error) {
    return "C", nil
  },
)

combinedPromise := promise.All(promiseA, promiseB)
result := combinedPromise.GetResult()
if result.Error != nil {
  // handle error
}
values := result.Value.([]interface{})

promiseAValue = values[0]
// promiseAValue == "A"
promiseBValues = values[1].([]interface{})
// promiseBValues[0] == "B"
// promiseBValues[1] == "C"
```





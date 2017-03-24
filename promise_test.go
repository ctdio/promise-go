package promise

import (
  "time"
  "testing"
)

type testStruct struct {
  String string
  Integer int
}

func (t testStruct) ReturnValues () (string, int) {
  return t.String, t.Integer
}

/**
 * Test to ensure that promise returns a Result struct
 * that can be consumed
 */
func TestNewPromiseResult (t *testing.T) {
  strVal := "yay promises"
  intVal := 100

  promise := New(func () (interface{}, error) {
    var value testStruct

    value.String = strVal
    value.Integer = intVal

    return value, nil
  })

  result, ok := <-promise
  if ok != true {
    t.Fatal("result from promise was not received from the channel")
  }

  newVal, isTestStruct := result.Value.(testStruct)

  // ensure returned value can be type casted
  // back to a usable object for the consumer
  if isTestStruct != true {
    t.Fatal("Failed")
  }

  val1, val2 := newVal.ReturnValues()
  if val1 != strVal || val2 != intVal {
    t.Fail()
  }
}

// make sure the promise is truly async
func TestNewPromiseIsAsync (t *testing.T) {
  promiseComplete := false

  promise := New(func () (interface{}, error) {
    time.Sleep(1 * time.Second)
    promiseComplete = true

    return promiseComplete, nil
  })

  if promiseComplete == true {
    t.Fatal("function passed to promise is not running asynchronously")
  }

  result, ok := <-promise
  if ok != true {
    t.Fatal("result from promise was not received from the channel")
  }

  if result.Value != true || promiseComplete != true {
    t.Fail()
  }
}

func TestPromiseChannelCloses (t *testing.T) {
  promise := New(func () (interface{}, error) {
    return true, nil
  })

  result, ok := <-promise
  if result.Value != true && ok != true {
    t.Fatal("result from promise was not received from the channel")
  }

  result, stillOpen := <-promise

  if result != nil || stillOpen == true {
    t.Fatal("Channel not closed")
  }
}

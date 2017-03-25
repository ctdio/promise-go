package promise

import (
  "time"
  "errors"
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
func TestCreatePromiseResult (t *testing.T) {
  strVal := "yay promises"
  intVal := 100

  promise := Create(func () (interface{}, error) {
    var value testStruct

    value.String = strVal
    value.Integer = intVal

    return value, nil
  })

  result := promise.GetResult()
  if result.Error != nil {
    t.Fatal("Unexpected error returned from promise")
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

/**
 * Test to ensure that promise returns a Result struct
 * that can be consumed
 */
func TestCreatePromiseError (t *testing.T) {
  promise := Create(func () (interface{}, error) {
    var value testStruct
    return value, errors.New("Error")
  })

  result := promise.GetResult()
  if result.Error == nil {
    t.Fatal("Error should have been returned from promise")
  }
}

/**
 * Test to ensure that an already resolved promise just returns
 * the same result
 */
func TestPromiseGetResult (t *testing.T) {
  promise := Create(func () (interface{}, error) {
    var value testStruct
    return value, nil
  })

  result := promise.GetResult()
  if result.Error != nil {
    t.Fatal("Unexpected error returned from promise")
  }
  if promise.settled != true {
    t.Fatal("Promise should have been marked as settled")
  }
  if result != promise.GetResult() {
    t.Fatal("Promise should have returned the same result")
  }
}

/**
 * Test to ensure that an already resolved promise just returns
 * the same result
 */
func TestPromiseChannelError (t *testing.T) {
  promise := Create(func () (interface{}, error) {
    time.Sleep(1 * time.Second)
    return 2, nil
  })

  close(promise.channel)

  result := promise.GetResult()
  if result.Error == nil {
    t.Fatal("Error should have been returned from promise")
  }
}

// make sure the promise is truly async
func TestCreatePromiseIsAsync (t *testing.T) {
  promiseComplete := false

  asyncFunc := func () (interface{}, error) {
    time.Sleep(1 * time.Second)
    promiseComplete = true

    return promiseComplete, nil
  }

  promise := Create(asyncFunc)

  if promiseComplete == true {
    t.Fatal("Function passed to promise is not running asynchronously")
  }

  result := promise.GetResult()

  if result.Value != true || promiseComplete != true {
    t.Fail()
  }
}

func TestPromiseAll (t *testing.T) {
  promiseA := Create(func () (interface{}, error) {
    time.Sleep(2 * time.Second)
    return true, nil
  })
  promiseB := Create(func () (interface{}, error) {
    time.Sleep(1 * time.Second)
    return false, nil
  })

  combined := All(promiseA, promiseB)

  result := combined.GetResult()
  if result.Error != nil {
    t.Fatal("Unexpected error from promises")
  }
  values := result.Value.([]interface{})
  resultLength := len(values)
  if resultLength != 2 {
    t.Fatalf("Expected result length to equal 2, got %s instead", resultLength)
  }

  if values[0] != true || values[1] != false {
    t.Fatal("Values in the result did not match the promised values")
  }
}

/**
 * Ensure that channels are closed after a result is returned
 * and that promises have the expected values
 */
func TestPromiseCleanup (t *testing.T) {
  promise := Create(func () (interface{}, error) {
    return true, nil
  })

  result := promise.GetResult()
  if result.Error != nil {
    t.Fatal("Unexpected error returned from promise")
  }
  if result.Value != true {
    t.Fatal("Invalid result received from promise")
  }

  if _, stillOpen := <-promise.channel; stillOpen {
    t.Fatal("Channel not closed")
  }

  if !promise.settled {
    t.Fatal("Result was not marked as settled result was retrieved")
  }

  if result != promise.result {
    t.Fatal("The internally stored result value does match was was returned from GetResult")
  }
}

func TestPromiseCreateAll (t *testing.T) {
  valueA := "value a"
  valueB := "value b"
  promise := CreateAll(
    func () (interface{}, error) {
      time.Sleep(1 * time.Second)
      return valueA, nil
    },
    func () (interface{}, error) {
      time.Sleep(2 * time.Second)
      return valueB, nil
    },
  )

  result := promise.GetResult()
  if result.Error != nil {
    t.Fatal("Unexpected error returned from promise")
  }

  values := result.Value.([]interface{})
  if len(values) != 2 {
    t.Fatal("Returned result does not match number of promises passed in")
  }

  if values[0] != valueA || values[1] != valueB {
    t.Fatal("Returned values do not match what was returned from async functions")
  }
}

func TestPromiseCreateAllError (t *testing.T) {
  valueA := "value a"
  valueB := "value b"

  promise := CreateAll(
    func () (interface{}, error) {
      time.Sleep(1 * time.Second)
      return valueA, errors.New("This should be part of the result")
    },
    func () (interface{}, error) {
      time.Sleep(2 * time.Second)
      return valueB, nil
    },
  )

  result := promise.GetResult()
  if result.Error == nil {
    t.Fatal("Error should have been returned from promise")
  }
}

func TestSingleAndCombinedPromises (t *testing.T) {
  valueA := "a"
  valueB := "b"
  valueC := "c"

  singlePromise := Create(func () (interface{}, error) {
    return valueA, nil
  })

  multiPromises := CreateAll(
    func () (interface{}, error) {
      time.Sleep(1 * time.Second)
      return valueB, nil
    },
    func () (interface{}, error) {
      time.Sleep(2 * time.Second)
      return valueC, nil
    },
  )

  // await all
  combinedPromises := All(singlePromise, multiPromises)
  result := combinedPromises.GetResult()
  if result.Error != nil {
    t.Fatal("Error retrieving value from promise")
  }

  values := result.Value.([]interface{})
  if len(values) != 2 {
    t.Fatal("Number of values do not match")
  }

  if values[0] != valueA {
    t.Fatal("Expected first value to match return value from singlePromise")
  }

  secondSet := values[1].([]interface{})
  if len(secondSet) != 2 {
    t.Fatal("Number of values do not match")
  }
  if secondSet[0] != valueB {
    t.Fatal("Expected secondset's first value to match return value from the first promise set")
  }
  if secondSet[1] != valueC {
    t.Fatal("Expected secondset's second value to match return value from the second promise set")
  }
}

package promise

import "errors"

type Result struct {
  Value interface{}
  Error error
  index int
}

type Promise struct {
  channel chan *Result
  settled bool
  result *Result
}

func (p *Promise) GetResult () *Result {
  // promise was already settled
  if p.settled {
    return p.result
  }

  var result *Result = nil

  // otherwise wait for result
  result, ok := <-p.channel
  if !ok {
    result = &Result{}
    result.Error = errors.New("Channel unexpectedly closed")
  }
  p.settled = true
  p.result = result // store result if we have one

  return result
}

type promiseFunction func () (interface{}, error)

/**
 * Creates a simple channel that returns the result
 * of the given function
 *
 * The value that is published is that of the Result struct
 */
func Create (function promiseFunction) *Promise {
  channel := make(chan *Result)

  // run given function in a goroutine
  go func () {
    defer func () {
      recover()
    }()
    res, err := function()
    channel <- &Result{res, err, 0}
    close(channel)
  }()

  return &Promise{channel, false, nil}
}

/**
 * awaits for results to be returned from the aggregate channel
 */
func awaitAggregatedResults (aggregateChan chan *Result, expectedResults int) ([]interface{}, error) {
  resultSlice := make([]interface{}, expectedResults)
  count := 0
  // read result from the aggregate channel
  for result := range aggregateChan {
    // if there was an error, close the channel
    // and return immediately
    if result.Error != nil {
      return nil, result.Error
    }
    resultSlice[result.index] = result.Value

    // if the count matches to total functions passed in, stop listening to channel
    // and return the result
    if count++; count == expectedResults {
      break
    }
  }

  return resultSlice, nil
}

/**
 * Waits for all promises to return a value (or fail)
 *
 * If any failure happens, this method will return the error immediately
 */
func All (promises ...*Promise) (*Promise) {
  awaitAllPromises := func () ([]interface{}, error) {
    promiseCount := len(promises)
    aggregateChan := make(chan *Result)

    for i, v := range promises {
      // copy values into local vars
      index := i
      promise := v

      // push result when ready
      go func () {
        defer func () {
          recover()
        }()
        result := promise.GetResult()
        result.index = index
        aggregateChan <- result
      }()
    }

    result, err := awaitAggregatedResults(aggregateChan, promiseCount)
    close(aggregateChan)

    return result, err
  }

  channel := make(chan *Result)

  // create a goroutine for awaiting all results
  go func () {
    res, err := awaitAllPromises()
    channel <- &Result{res, err, 0}
    close(channel)
  }()

  return &Promise{channel, false, nil}
}

/**
 * Invokes all of the functions and returns a *Promise
 * which can be used to block when the consumer is ready
 *
 */
func CreateAll (functions ...promiseFunction) *Promise {
  functionCount := len(functions)
  promises := make([]*Promise, functionCount)

  // for each function, start a goroutine
  for i, function := range functions {
    promises[i] = Create(function)
  }

  return All(promises...)
}


package promise

import "errors"

type Result struct {
  Value interface{}
  Error error
  index int
}

type CombinedResult struct {
  Values []interface{}
  Error error
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

  // otherwise wait for result
  result, ok := <-p.channel
  if !ok {
    result.Error = errors.New("Channel unexpectedly closed")
  }
  p.settled = true
  p.result = result // store result if we have one

  return result
}

type promiseFunction func () (interface{}, error)

type CombinedPromise struct {
  channel chan *CombinedResult
}

func (p *CombinedPromise) GetResult () *CombinedResult {
  result, ok := <-p.channel
  if !ok {
    result.Error = errors.New("Channel unexpectedly closed")
  }
  return result
}

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
    res, err := function()
    channel <- &Result{res, err, 0}
    close(channel)
  }()

  return &Promise{channel, false, nil}
}

/**
 * awaits for results to be returned from the aggregate channel
 */
func awaitAggregateResults (aggregateChan chan *Result, expectedResults int) ([]interface{}, error) {
  resultSlice := make([]interface{}, expectedResults)
  count := 0
  // read result from the aggregate channel
  for result := range aggregateChan {
    // if there was an error, close the channel
    // and return immediately
    if result.Error != nil {
      close(aggregateChan)
      return nil, result.Error
    }
    resultSlice[result.index] = result.Value

    // if the count matches to total functions passed in,
    // close the channel
    if count++; count == expectedResults {
      close(aggregateChan)
    }
  }

  return resultSlice, nil
}

/**
 * Waits for all promises to return a value (or fail)
 *
 * If any failure happens, this method will return the error immediately
 */
func All (promises ...*Promise) (*CombinedPromise) {
  awaitAllPromises := func () ([]interface{}, error) {
    promiseCount := len(promises)
    aggregateChan := make(chan *Result)

    for i, v := range promises {
      // copy values into local vars
      index := i
      promise := v

      // push result when ready
      go func () {
        result := promise.GetResult()
        result.index = index
        aggregateChan <- result
      }()
    }

    return awaitAggregateResults(aggregateChan, promiseCount)
  }

  channel := make(chan *CombinedResult)

  // create a goroutine for awaiting all results
  go func () {
    res, err := awaitAllPromises()
    channel <- &CombinedResult{res, err}
    close(channel)
  }()

  return &CombinedPromise{channel}
}

/**
 * Invokes all of the functions and returns a *CombinedPromise
 * which can be used to block when the consumer is ready
 *
 */
func CreateAll (functions ...promiseFunction) *CombinedPromise {
  functionCount := len(functions)
  promises := make([]*Promise, functionCount)

  // for each function, start a goroutine
  for i, function := range functions {
    promises[i] = Create(function)
  }

  return All(promises...)
}


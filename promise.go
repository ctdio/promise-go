package promise

type Result struct {
  Value interface{}
  Error error
}

type asyncFunc func () (interface{}, error)

/**
 * Creates a simple channel that returns the result
 * of the given function
 *
 * The value that is published is that of the Result struct
 */
func New (function asyncFunc) chan *Result {
  channel := make(chan *Result)

  // run given function in a goroutine
  go func () {
    res, err := function()
    channel <- &Result{res, err}
    close(channel)
  }()

  return channel
}

package asyncaller

type AsyncFunc[T any] func() (T, error)

type AsyncResult[T any] struct {
	result T
	err    error
	done   chan struct{}
}

func NewAsyncResult[K any]() *AsyncResult[K] {
	return &AsyncResult[K]{
		done: make(chan struct{}),
	}
}

func (a *AsyncResult[T]) Get() (T, error) {
	<-a.done
	return a.result, a.err
}

func (a *AsyncResult[T]) IsDone() bool {
	select {
	case <-a.done:
		return true
	default:
		return false
	}
}

func Exec[S any](fn AsyncFunc[S]) *AsyncResult[S] {
	result := NewAsyncResult[S]()

	go func() {
		r, err := fn()
		result.result = r
		result.err = err
		close(result.done)
	}()

	return result
}

package flow

import (
    "context"
)

var _ interface {
    Pushable[any]
    Closable[any]
    Listenable[any]
} = (*Stream[any])(nil)

// Stream is a Flow implementation that emits ordered values until closed.
type Stream[T any] struct {
    base[T]
}

// Close finalizes the flow by pushing the provided values and marking it as closed.
func (s *Stream[T]) Close(value ...T) {
    s.init()
    s.write(func() {
        s.push(value)
        s.close()
    })
}

// Push adds the provided values into the flow.
func (s *Stream[T]) Push(value ...T) {
    s.init()
    s.write(func() { s.push(value) })
}

// Listen returns a channel that emits values from the head of the flow.
// The channel closes when the context is canceled or the flow is closed.
func (s *Stream[T]) Listen(ctx context.Context) (c <-chan T) {
    s.init()
    s.read(func() { c = s.head.listen(ctx) })
    return
}

func (s *Stream[T]) init() {
    s.base.init(nil)
}

// push adds the provided values into the flow and resets its state.
func (s *Stream[T]) push(values []T) {
    s.base.push(values)
    s.reset()
}

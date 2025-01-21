package flow

import "context"

type Listenable[T any] interface {
    // Listen returns a channel that emits values pushed into the flow.
    // The channel closes when either the Closable.Close method is called or the provided context is cancelled.
    // The values received by the channel will be in the same order they were pushed.
    Listen(ctx context.Context) <-chan T
}

type Pushable[T any] interface {
    // Push allows multiple values to be pushed into the flow.
    // This operation is thread-safe and does not block the caller.
    Push(...T)
}

type Closable[T any] interface {
    // Close pushes final values (if any) and marks the flow as closed.
    // After Close is called, all channels produced by Listenable.Listen will close after emitting the final values.
    // Once closed, Resettable.Reset will no longer have any effect.
    // Close can be safely called multiple times, but only the first invocation will push values and close the flow.
    Close(...T)
}

type Resettable[T any] interface {
    // Reset clears the flow and pushes new values into it.
    // If the flow is already closed (via Closable.Close), Reset will have no effect.
    Reset(...T)
}

// Flow is an interface that combines [Listenable], [Pushable], [Closable], and [Resettable].
type Flow[T any] interface {
    Listenable[T]
    Pushable[T]
    Closable[T]
    Resettable[T]
}

package flow

import (
    "context"
)

// Never represents a Flow that never emits any value and closes only when the Listen context is cancelled.
type Never[T any] struct{}

func (Never[T]) Listen(ctx context.Context) <-chan T {
    c := make(chan T)
    go func() {
        defer close(c)
        <-ctx.Done()
    }()
    return c
}

// Empty represents a Flow that immediately closes without emitting any values.
type Empty[T any] struct{}

func (Empty[T]) Listen(context.Context) <-chan T {
    c := make(chan T)
    close(c)
    return c
}

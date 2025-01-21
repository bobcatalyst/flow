package flow

import (
    "context"
)

var _ Flow[any] = (*Repeater[any])(nil)

// Repeater is a Flow implementation that repeats values until reset.
type Repeater[T any] struct {
    base[T]
    // ref keeps track of the start node use when calling Listen.
    // This will allow all pushed values to be kept and repeated until reset.
    ref *node[T]
}

// Close finalizes the flow by pushing the provided values and marking it as closed.
func (r *Repeater[T]) Close(value ...T) {
    r.init()
    r.write(func() {
        r.push(value)
        r.close()
    })
}

// Reset clears the flow's current state and initializes it with the provided values.
func (r *Repeater[T]) Reset(value ...T) {
    r.init()
    r.write(func() {
        r.reset()
        r.push(value)
        if !r.head.final {
            // If the flow is not closed, update the reference node.
            r.ref = r.head
        }
    })
}

// Push adds the provided values into the flow.
func (r *Repeater[T]) Push(value ...T) {
    r.init()
    r.write(func() { r.push(value) })
}

// Listen returns a channel that emits values starting from the current reference node.
// The channel will close when the context is canceled or the flow is closed.
func (r *Repeater[T]) Listen(ctx context.Context) (c <-chan T) {
    r.init()
    r.read(func() { c = r.ref.listen(ctx) })
    return
}

// init initializes the flow, ensuring the base is set up and the reference node points to the start of the buffer.
func (r *Repeater[T]) init() {
    // Set the reference node to the head of the flow to buffer values.
    r.base.init(func() { r.ref = r.head })
}

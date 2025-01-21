package flow

import "context"

// node represents a single element in a linked list structure used for managing values in the flow.
type node[T any] struct {
    wait     chan struct{} // Channel to signal when a new node is available.
    next     *node[T]      // Next node in the sequence.
    hasValue bool          // Indicates whether this node holds a value.
    value    T             // The value held by this node.
    final    bool          // Marks whether this node is the final one in the flow.
}

// newNode creates a new node ready to be waited on.
func newNode[T any]() *node[T] {
    return &node[T]{
        wait: make(chan struct{}),
    }
}

// push creates a new node with the provided value and links it to the current node.
func (n *node[T]) push(v T) *node[T] {
    return n.setNext(func() *node[T] {
        nn := newNode[T]()
        nn.hasValue = true
        nn.value = v
        return nn
    })
}

// reset clears the flow by creating a new empty node and linking it to the current node.
func (n *node[T]) reset() *node[T] {
    return n.setNext(newNode)
}

// close marks the flow as complete by creating a final node and linking it to the current node.
func (n *node[T]) close() *node[T] {
    return n.setNext(func() *node[T] {
        nn := newNode[T]()
        nn.final = true
        return nn
    })
}

// setNext sets the next node in the sequence using a constructor to generate the new node.
// If the current node is already final, it does nothing. If a next node already exists,
// the method recursively delegates to the next node.
func (n *node[T]) setNext(fn func() *node[T]) *node[T] {
    if n.final {
        return n
    }

    // make sure that we're operating on the last node
    for ; n.next != nil; n = n.next {
        if n.final {
            return n
        }
    }

    nn := fn()
    n.next = nn
    close(n.wait)
    return nn
}

// listen starts listening to values emitted by this node and its subsequent nodes.
// Returns a channel that produces the values. The channel closes when the context is canceled
// or when a final node is reached.
func (n *node[T]) listen(ctx context.Context) <-chan T {
    c := make(chan T)
    if ctx == nil {
        ctx = context.Background()
    }
    go sendValues(ctx, n, c)
    return c
}

// sendValues traverses the linked list starting from the current node,
// sending values to the provided channel until the context is canceled or a final node is reached.
func sendValues[T any](ctx context.Context, n *node[T], c chan<- T) {
    defer close(c) // Ensure the channel is closed when the function exits.
    for ; ctx.Err() == nil && n != nil && !n.final; n = n.next {
        sendValue(n, ctx, c)
        nodeWait(n, ctx)
    }
}

// sendValue sends the value of the given node to the channel if the context is still active and the node has a value.
func sendValue[T any](n *node[T], ctx context.Context, c chan<- T) {
    if n.hasValue {
        select {
        case <-ctx.Done():
        case c <- n.value:
        }
    }
}

// nodeWait waits for the next node to be ready or for the context to be canceled.
func nodeWait[T any](n *node[T], ctx context.Context) {
    select {
    case <-ctx.Done():
    case <-n.wait:
    }
}

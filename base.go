package flow

import "sync"

// base is a generic struct that provides thread-safe operations for managing a linked list of nodes.
type base[T any] struct {
    lock sync.RWMutex
    once sync.Once
    head *node[T]
}

// init initializes the base structure with a new head node and an optional initialization function.
// Will only run the first time it is called.
func (b *base[T]) init(fn func()) {
    b.once.Do(func() {
        b.head = newNode[T]()
        if fn != nil {
            fn()
        }
    })
}

// write performs a thread-safe write operation by acquiring the write lock.
func (b *base[T]) write(fn func()) {
    withLock(&b.lock, fn)
}

// read performs a thread-safe read operation by acquiring the read lock.
func (b *base[T]) read(fn func()) {
    withLock(b.lock.RLocker(), fn)
}

// withLock is a utility function that locks a given locker, executes the function, and then always unlocks it.
func withLock(l sync.Locker, fn func()) {
    l.Lock()
    defer l.Unlock()
    fn()
}

// push adds multiple values to the linked list by pushing them into the head node sequentially.
func (b *base[T]) push(v []T) {
    for _, v := range v {
        b.head = b.head.push(v)
    }
}

// reset clears the flow by replacing the head with an empty node.
func (b *base[T]) reset() {
    b.head = b.head.reset()
}

// close closes the flow by replacing the head with a final node.
func (b *base[T]) close() {
    b.head = b.head.close()
}

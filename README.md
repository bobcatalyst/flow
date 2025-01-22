# Flow Library

## Overview

The **Flow** library provides a thread-safe implementation for managing asynchronous data flows. 
It includes several implementations, such as `Stream`, `Repeater`, `Never`, and `Empty`, each with different behaviors.

## Installation

Include the library in your Go project:

```bash
go get github.com/bobcatalyst/flow
```

## Implementations

### Stream
A `Stream` represents a continuous flow of values that can be closed, pushed to, and listened to.

#### Example Usage

```go
var s flow.Stream[int]
s.Push(1, 2, 3)
c := s.Listen(context.Background())
s.Push(4, 5, 6)
for val := range c {
    fmt.Println(val)
}
```

##### Output:

```
4
5
6
```

### Repeater
A `Repeater` repeats past values.

#### Example Usage

```go
var r flow.Repeater[int]
r.Push(1, 2, 3)
c := s.Listen(context.Background())
r.Push(4, 5, 6)
for val := range c {
    fmt.Println(val)
}
```

##### Output:

```
1
2
3
4
5
6
```

### Never
A `Never` flow never emits any values and only closes when the context is canceled.

#### Example Usage
```go
var n flow.Never[int]
c := n.Listen(context.Background())
// Will block until the context is canceled.
```

### Empty
An `Empty` flow immediately closes without emitting any values.

#### Example Usage
```go
var e flow.Empty[int]
c := e.Listen(context.Background())
// Channel `c` is already closed.
```

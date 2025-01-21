package flow

import (
    "context"
    "testing"
    "time"
)

func TestEmpty_Listen(t *testing.T) {
    var e Empty[any]
    c := e.Listen(context.Background())
    select {
    case <-c:
    case <-time.After(time.Second):
        t.Fail()
    }
}

func TestNever_Listen(t *testing.T) {
    var e Never[any]
    ctx, cancel := context.WithCancel(context.Background())
    c := e.Listen(ctx)
    select {
    case <-c:
        t.Fail()
    case <-time.After(time.Second):
    }
    cancel()
    select {
    case <-c:
    case <-time.After(time.Second):
        t.Fail()
    }
}

package diodes

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type hookContext struct {
	context.Context
	onDone func()
}

func (h *hookContext) Done() <-chan struct{} {
	if h.onDone != nil {
		h.onDone()
	}
	return h.Context.Done()
}

type hookDiode struct {
	Diode
	onEmptyTryNext func()
}

func (h *hookDiode) TryNext() (GenericDataType, bool) {
	data, ok := h.Diode.TryNext()
	if !ok && h.onEmptyTryNext != nil {
		h.onEmptyTryNext()
	}
	return data, ok
}

func TestWaiterSetDoesNotLoseWakeup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var emptySeen atomic.Bool
	var once sync.Once
	var w *Waiter

	payload := "hello"
	setDone := make(chan struct{})

	hd := &hookDiode{
		Diode: NewManyToOne(4, AlertFunc(func(int) {})),
		onEmptyTryNext: func() {
			emptySeen.Store(true)
		},
	}

	hctx := &hookContext{
		Context: ctx,
		onDone: func() {
			// Waiter.Next calls w.isDone() -> w.ctx.Done() while holding w.mu,
			// immediately after TryNext() returned empty and before w.c.Wait().
			if !emptySeen.Load() {
				return
			}
			once.Do(func() {
				go func() {
					w.Set(GenericDataType(&payload))
					close(setDone)
				}()
				// If Waiter.Set does not lock w.mu before w.c.Broadcast(), Set
				// completes and broadcasts here before Next() enters w.c.Wait(),
				// causing a permanent lost wakeup. With the lock held in Set,
				// Set blocks on w.mu until Next() parks in w.c.Wait().
				select {
				case <-setDone:
				case <-time.After(25 * time.Millisecond):
				}
			})
		},
	}

	w = NewWaiter(hd, WithWaiterContext(hctx))

	gotCh := make(chan GenericDataType, 1)
	go func() {
		gotCh <- w.Next()
	}()

	select {
	case got := <-gotCh:
		if got == nil {
			t.Fatal("Next() returned nil instead of buffered payload")
		}
		if *(*string)(got) != payload {
			t.Fatalf("Next() = %q, want %q", *(*string)(got), payload)
		}
	case <-time.After(500 * time.Millisecond):
		cancel()
		<-gotCh
		t.Fatal("Waiter.Next() hung due to lost wakeup in Waiter.Set()")
	}

	select {
	case <-setDone:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Waiter.Set() did not complete")
	}
}

func TestWaiterContextCancelWakesNext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	w := NewWaiter(NewManyToOne(4, AlertFunc(func(int) {})), WithWaiterContext(ctx))

	gotCh := make(chan GenericDataType, 1)
	go func() {
		gotCh <- w.Next()
	}()

	cancel()

	select {
	case got := <-gotCh:
		if got != nil {
			t.Fatalf("Next() after cancel = %v, want nil", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Waiter.Next() did not wake up after context cancellation")
	}
}

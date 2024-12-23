package ctxlock

import (
	"context"
)

type Lock struct {
	c chan struct{}
}

func New() Lock {
	l := Lock{
		c: make(chan struct{}, 1),
	}
	l.c <- struct{}{}
	return l
}

func (l *Lock) Lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.c:
		return nil
	}
}

func (l *Lock) Unlock() {
	select {
	case l.c <- struct{}{}:
	default:
		panic("unlock of unlocked lock")
	}
}

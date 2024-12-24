package keylock

import (
	"context"
	"fmt"

	"github.com/c4t-but-s4d/fastad/pkg/ctxlock"
)

type Lock[T comparable] struct {
	mu ctxlock.Lock

	locks map[T]ctxlock.Lock
}

func New[T comparable](prealloc ...T) Lock[T] {
	l := Lock[T]{
		locks: make(map[T]ctxlock.Lock, len(prealloc)),
	}
	for _, t := range prealloc {
		l.locks[t] = ctxlock.New()
	}
	return l
}

func (l *Lock[T]) Lock(ctx context.Context, key T) (func(), error) {
	if err := l.mu.Lock(ctx); err != nil {
		return nil, fmt.Errorf("acquiring global lock: %w", err)
	}
	defer l.mu.Unlock()

	cl, ok := l.locks[key]
	if !ok {
		cl = ctxlock.New()
		l.locks[key] = cl
	}
	if err := cl.Lock(ctx); err != nil {
		return nil, fmt.Errorf("acquiring inner lock: %w", err)
	}
	return func() {
		cl.Unlock()
	}, nil
}

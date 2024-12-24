package keylock

import (
	"context"
	"fmt"
	"slices"

	"golang.org/x/exp/constraints"

	"github.com/c4t-but-s4d/fastad/pkg/ctxlock"
)

type MultiLock[T constraints.Ordered] struct {
	mu ctxlock.Lock

	locks map[T]ctxlock.Lock
}

func NewMulti[T constraints.Ordered](prealloc ...T) MultiLock[T] {
	l := MultiLock[T]{
		locks: make(map[T]ctxlock.Lock, len(prealloc)),
	}
	for _, t := range prealloc {
		l.locks[t] = ctxlock.New()
	}
	return l
}

func (l *MultiLock[T]) Lock(ctx context.Context, keys ...T) (func(), error) {
	slices.Sort(keys)

	if err := l.mu.Lock(ctx); err != nil {
		return nil, fmt.Errorf("acquiring global lock: %w", err)
	}
	defer l.mu.Unlock()

	locked := make([]ctxlock.Lock, 0, len(keys))
	cleanup := func() {
		for _, cl := range locked {
			cl.Unlock()
		}
	}

	for _, key := range keys {
		cl, ok := l.locks[key]
		if !ok {
			cl = ctxlock.New()
			l.locks[key] = cl
		}
		if err := cl.Lock(ctx); err != nil {
			cleanup()
			return nil, fmt.Errorf("acquiring inner lock: %w", err)
		}
		locked = append(locked, cl)
	}

	return cleanup, nil
}

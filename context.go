// Package trace provides a wrapper around OpenTelemetry's trace API
// with support for goroutine-bound context propagation.
package trace

import (
	"context"
	"sync"
)

var (
	contexts   = make(map[uint64]*contextWrapper)
	mxContexts sync.RWMutex
)

// contextWrapper is a wrapper around a context.Context that tracks its parent context.
type contextWrapper struct {
	context.Context
	id     uint64
	parent *contextWrapper
}

func enterContext(ctx context.Context) *contextWrapper {
	id := curGoroutineID()
	mxContexts.Lock()
	defer mxContexts.Unlock()
	parentOrNil := contexts[id]
	next := &contextWrapper{
		Context: ctx,
		id:      id,
		parent:  parentOrNil,
	}
	contexts[id] = next

	return next
}

func (c *contextWrapper) exit() {
	id := c.id
	mxContexts.Lock()
	defer mxContexts.Unlock()
	current := contexts[id]
	if current == nil {
		return
	}
	if current.parent == nil {
		delete(contexts, id)
	} else {
		contexts[id] = current.parent
	}
}

// CurrentContext returns the current goroutine-local context, or
// context.Background() if there is no current context.
func CurrentContext() context.Context {
	id := curGoroutineID()
	mxContexts.RLock()
	defer mxContexts.RUnlock()
	current := contexts[id]
	if current == nil {
		return context.Background()
	}
	return current
}

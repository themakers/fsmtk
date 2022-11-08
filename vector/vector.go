package vector

import (
	"context"
	"sync"
)

type ProcessFunc func(ctx context.Context)
type SetVectorFunc func(vector bool)

type Vector struct {
	vector  bool
	working bool
	lock    sync.Mutex

	cancelCurrentWork context.CancelFunc

	rootCtx context.Context
	process ProcessFunc
}

func New(ctx context.Context, process ProcessFunc) *Vector {
	return &Vector{
		cancelCurrentWork: func() {},

		rootCtx: ctx,
		process: process,
	}
}

func (v *Vector) withLock(fn func()) {
	v.lock.Lock()
	defer v.lock.Unlock()
	fn()
}

func (v *Vector) Set(vec bool) {
	v.withLock(func() {
		if v.vector != vec {
			v.vector = vec
			v.trigger()
		}
	})
}

func (v *Vector) trigger() {

	if v.vector != v.working {
		if v.vector { //> Gonna run

			ctx, cancel := context.WithCancel(v.rootCtx)
			v.cancelCurrentWork = cancel
			go v.run(ctx)

		} else { //> Gonna die

			v.cancelCurrentWork()

		}
	}

}

func (v *Vector) run(ctx context.Context) {
	defer v.withLock(func() {
		v.working = false
		v.trigger()
	})

	v.withLock(func() {
		v.working = true
		v.trigger()
	})

	v.process(ctx)
}

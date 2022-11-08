package transform

import (
	"reflect"
	"sync"
)

type TransformFunc[Primary, Derived any] func(Primary) Derived

func Dummy[PD any](pd PD) PD {
	return pd
}

type DerivedChangedFunc[Derived any] func(Derived)

type Transform[Primary, Derived any] struct {
	prim Primary
	derv Derived

	lock sync.Mutex

	transformFn TransformFunc[Primary, Derived]
	outFn       DerivedChangedFunc[Derived]
}

func New[Primary, Derived any](
	transformFn TransformFunc[Primary, Derived],
	outFn DerivedChangedFunc[Derived],
) *Transform[Primary, Derived] {

	return &Transform[Primary, Derived]{
		transformFn: transformFn,
		outFn:       outFn,
	}
}

func (t *Transform[Primary, Derived]) withLock(fn func()) {
	t.lock.Lock()
	defer t.lock.Unlock()
	fn()
}

func (t *Transform[Primary, Derived]) set(prim Primary) {
	t.prim = prim
	derv := t.transformFn(prim)

	if !reflect.DeepEqual(t.derv, derv) {
		t.derv = derv
		t.outFn(t.derv)
	}
}

func (t *Transform[Primary, Derived]) Set(prim Primary) {
	t.withLock(func() {
		t.set(prim)
	})
}

func (t *Transform[Primary, Derived]) Mutate(mfn func(Primary) Primary) {
	t.withLock(func() {
		prim := mfn(t.prim)
		t.set(prim)
	})
}

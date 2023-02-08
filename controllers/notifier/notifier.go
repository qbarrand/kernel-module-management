package notifier

import (
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/event"
)

type Notifier[T any] struct {
	ch  chan event.GenericEvent
	m   sync.RWMutex
	val T
}

func New[T any](val T) *Notifier[T] {
	return &Notifier[T]{
		ch:  make(chan event.GenericEvent),
		val: val,
	}
}

func (sc *Notifier[T]) Get() T {
	sc.m.RLock()
	defer sc.m.RUnlock()

	return sc.val
}

func (sc *Notifier[T]) Set(val T) {
	sc.m.Lock()
	sc.val = val
	defer sc.m.Unlock()

	sc.ch <- event.GenericEvent{}
}

func (sc *Notifier[T]) EventChannel() <-chan event.GenericEvent {
	return sc.ch
}

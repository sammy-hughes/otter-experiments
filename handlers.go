package cache

import (
	"time"

	"github.com/maypok86/otter"
)

type ShiftTTLFunc[K comparable, V any] func(K, V, time.Duration) time.Duration
type BeforeLoadFunc[K comparable, V any] func(K)
type AfterLoadFunc[K comparable, V any] func(K)
type BeforeStoreFunc[K comparable, V any] func(K)
type AfterStoreFunc[K comparable, V any] func(K, V, bool)
type BeforeDeleteFunc[K comparable, V any] func(K)
type AfterDeleteFunc[K comparable, V any] func(K, V, otter.DeletionCause)
type OnStartFunc[K comparable, V any] func(*otter.CacheWithVariableTTL[K, V]) error

func DefaultAfterStart[K comparable, V any](_ *otter.Cache[K, V]) error {
	return nil
}

type Handlers[K comparable, V any] struct {
	ShiftTTL             ShiftTTLFunc[K, V]
	BeforeLoad           BeforeLoadFunc[K, V]
	AfterLoad            AfterLoadFunc[K, V]
	BeforeStore          BeforeStoreFunc[K, V]
	AfterStore           AfterStoreFunc[K, V]
	BeforeExplicitDelete BeforeDeleteFunc[K, V]
	AfterExplicitDelete  AfterDeleteFunc[K, V]
	OnStart              OnStartFunc[K, V]
}

func join1in[K comparable, V any, F ~func(K)](a, b F) F {
	switch {
	case a == nil:
		return b
	case b == nil:
		return a
	default:
		return func(k K) {
			a(k)
			b(k)
		}
	}
}

func join3in[K comparable, V any, X any, F ~func(K, V, X)](a, b F) F {
	switch {
	case a == nil:
		return b
	case b == nil:
		return a
	default:
		return func(k K, v V, x X) {
			a(k, v, x)
			b(k, v, x)
		}
	}
}

package cache

import (
	"time"

	"github.com/maypok86/otter"
)

type VisitorFunc[K comparable, V any] func(K, V) bool
type CostFunc[K comparable, V any] func(K, V) uint32

type Config[K comparable, V any] struct {
	Builder *otter.Builder[K, V]
	Handlers[K, V]
}

type Cache[K comparable, V any] struct {
	otter.CacheWithVariableTTL[K, V]
	Handlers[K, V]
}

func New[K comparable, V any](maxCapacity int, options ...CacheOption[K, V]) (Cache[K, V], error) {
	var result Cache[K, V]
	builder, err := otter.NewBuilder[K, V](maxCapacity)
	config := Config[K, V]{Builder: builder}

	if err != nil {
		return Cache[K, V]{}, err
	}

	for i := range options {
		err = options[i](&config)
		if err != nil {
			return Cache[K, V]{}, err
		}
	}

	result.CacheWithVariableTTL, err = builder.WithVariableTTL().Build()
	if err != nil {
		return Cache[K, V]{}, err
	}

	if config.Handlers.OnStart != nil {
		err = config.OnStart(&result.CacheWithVariableTTL)
		if err != nil {
			result.CacheWithVariableTTL = otter.CacheWithVariableTTL[K, V]{}
		}
	}

	return result, err
}

func (cache *Cache[K, V]) Delete(key K) {
	cache.Handlers.BeforeExplicitDelete(key)

	v, ok := cache.Get(key)
	if !ok {
		return
	}

	cache.Handlers.AfterExplicitDelete(key, v, otter.Explicit)
	cache.CacheWithVariableTTL.Delete(key)
}

func (cache *Cache[K, V]) DeleteByFunc(visitor VisitorFunc[K, V]) {
	cache.CacheWithVariableTTL.DeleteByFunc(visitor)
}

func (cache *Cache[K, V]) Has(key K) bool {
	cache.Handlers.BeforeLoad(key)

	ok := cache.CacheWithVariableTTL.Has(key)

	cache.Handlers.AfterLoad(key)
	return ok
}

func (cache *Cache[K, V]) Get(key K) (V, bool) {
	cache.Handlers.BeforeLoad(key)

	v, ok := cache.CacheWithVariableTTL.Get(key)
	if !ok {
		return v, false
	}

	cache.Handlers.AfterLoad(key)
	return v, ok
}

func (cache *Cache[K, V]) Range(visitor VisitorFunc[K, V]) {
	cache.CacheWithVariableTTL.Range(visitor)
}

func (cache *Cache[K, V]) Set(key K, value V, ttl time.Duration) bool {
	cache.Handlers.BeforeStore(key)

	ok := cache.CacheWithVariableTTL.Set(key, value, ttl)

	cache.Handlers.AfterStore(key, value, ok)
	return ok
}

func (cache *Cache[K, V]) SetIfAbsent(key K, value V, ttl time.Duration) bool {
	cache.Handlers.BeforeStore(key)

	ok := cache.CacheWithVariableTTL.SetIfAbsent(key, value, ttl)

	cache.Handlers.AfterStore(key, value, ok)
	return ok
}

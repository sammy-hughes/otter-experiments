package cache

import (
	"sync"

	"github.com/maypok86/otter"
)

type CacheOptionGroup[K comparable, V any] []CacheOption[K, V]

func (group CacheOptionGroup[K, V]) Apply(config *Config[K, V]) error {
	for i := range group {
		err := group[i](config)
		if err != nil {
			return err
		}
	}

	return nil
}

func WithMutex[K comparable, V any]() CacheOptionGroup[K, V] {
	mutex := sync.RWMutex{}

	return CacheOptionGroup[K, V]{
		WithBeforeLoadHandler[K, V](func(_ K) {
			mutex.RLock()
		}),
		WithAfterLoadHandler[K, V](func(_ K) {
			mutex.RUnlock()
		}),
		WithBeforeStoreHandler[K, V](func(_ K) {
			mutex.Lock()
		}),
		WithAfterStoreHandler(func(_ K, _ V, _ bool) {
			mutex.Unlock()
		}),
		WithBeforeExplicitDeleteHandler[K, V](func(_ K) {
			mutex.Lock()
		}),
		WithAfterExplicitDeleteHandler(func(_ K, _ V, _ otter.DeletionCause) {
			mutex.Unlock()
		}),
	}
}

func WithShardedMutex[K comparable, V any](sharder func(K) K, seedShardKeys ...K) CacheOptionGroup[K, V] {
	mutices := make(map[K]*sync.RWMutex)
	for i := range seedShardKeys {
		mutices[sharder(seedShardKeys[i])] = &sync.RWMutex{}
	}

	return CacheOptionGroup[K, V]{
		WithBeforeLoadHandler[K, V](func(k K) {
			mutex, ok := mutices[sharder(k)]
			if !ok {
				mutex = &sync.RWMutex{}
				mutices[sharder(k)] = mutex
			}

			mutex.RLock()
		}),
		WithAfterLoadHandler[K, V](func(k K) {
			mutex, ok := mutices[sharder(k)]
			if ok {
				// only unlock if shard exists already
				mutex.RUnlock()
			}
		}),
		WithBeforeStoreHandler[K, V](func(k K) {
			mutex, ok := mutices[sharder(k)]
			if !ok {
				mutex = &sync.RWMutex{}
				mutices[sharder(k)] = mutex
			}

			mutex.Lock()
		}),
		WithAfterStoreHandler(func(k K, _ V, _ bool) {
			mutex, ok := mutices[sharder(k)]
			if ok {
				// only unlock if shard exists already
				mutex.Unlock()
			}
		}),
		WithBeforeExplicitDeleteHandler[K, V](func(k K) {
			mutex, ok := mutices[sharder(k)]
			if !ok {
				mutex = &sync.RWMutex{}
				mutices[sharder(k)] = mutex
			}

			mutex.Lock()
		}),
		WithAfterExplicitDeleteHandler(func(k K, _ V, _ otter.DeletionCause) {
			mutex, ok := mutices[sharder(k)]
			if ok {
				// only unlock if shard exists already
				mutex.Unlock()
			}
		}),
	}
}

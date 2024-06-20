package cache

import (
	"time"

	"github.com/maypok86/otter"
)

type CacheOption[K comparable, V any] func(builder *Config[K, V]) error

func WithStatistics[K comparable, V any](config *Config[K, V]) error {
	config.Builder.CollectStats()

	return nil
}

func WithInitialCapacity[K comparable, V any](initialCapacity int) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Builder.InitialCapacity(initialCapacity)

		return nil
	}
}

func WithDeletionListener[K comparable, V any](listener AfterDeleteFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Builder.DeletionListener(listener)

		return nil
	}
}

func WithCostEstimates[K comparable, V any](coster CostFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Builder.Cost(coster)

		return nil
	}
}

func WithShiftTTLHandler[K comparable, V any](shifter ShiftTTLFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		var head ShiftTTLFunc[K, V]
		switch config.Handlers.ShiftTTL {
		case nil:
			head = func(_ K, _ V, ttl time.Duration) time.Duration {
				return ttl
			}
		default:
			head = config.Handlers.ShiftTTL
		}

		config.Handlers.ShiftTTL = func(k K, v V, ttl time.Duration) time.Duration {
			return shifter(k, v, head(k, v, ttl))
		}

		return nil
	}
}

func WithBeforeLoadHandler[K comparable, V any](seeLoad BeforeLoadFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Handlers.BeforeLoad = join1in[K, V](config.Handlers.BeforeLoad, seeLoad)

		return nil
	}
}

func WithAfterLoadHandler[K comparable, V any](sawLoad AfterLoadFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Handlers.AfterLoad = join1in[K, V](config.Handlers.AfterLoad, sawLoad)

		return nil
	}
}

func WithBeforeStoreHandler[K comparable, V any](seeStore BeforeStoreFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Handlers.BeforeStore = join1in[K, V](config.Handlers.BeforeStore, seeStore)

		return nil
	}
}

func WithAfterStoreHandler[K comparable, V any](sawStore AfterStoreFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Handlers.AfterStore = join3in(config.Handlers.AfterStore, sawStore)

		return nil
	}
}

func WithBeforeExplicitDeleteHandler[K comparable, V any](seeDelete BeforeDeleteFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Handlers.BeforeExplicitDelete = join1in[K, V](config.Handlers.BeforeExplicitDelete, seeDelete)

		return nil
	}
}

func WithAfterExplicitDeleteHandler[K comparable, V any](sawDelete AfterDeleteFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		config.Handlers.AfterExplicitDelete = join3in(config.Handlers.AfterExplicitDelete, sawDelete)

		return nil
	}
}

func WithOnStartHandler[K comparable, V any](starting OnStartFunc[K, V]) CacheOption[K, V] {
	return func(config *Config[K, V]) error {
		var head OnStartFunc[K, V]
		switch config.Handlers.OnStart {
		case nil:
			head = func(_ *otter.CacheWithVariableTTL[K, V]) error {
				return nil
			}
		default:
			head = config.Handlers.OnStart
		}

		config.Handlers.OnStart = func(cache *otter.CacheWithVariableTTL[K, V]) error {
			err := head(cache)
			if err != nil {
				return err
			}

			return starting(cache)
		}

		return nil
	}
}

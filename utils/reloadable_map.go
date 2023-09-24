package utils

import (
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

type Hashable interface {
	constraints.Integer | ~string
}

type MapDataLoader[T any] func(args ...any) ([]T, error)
type KeyProvider[K Hashable, T any] func(element T) K

type ReloadableMap[K Hashable, T any] struct {
	data           sync.Map
	loader         MapDataLoader[T]
	keyProvider    KeyProvider[K, T]
	reloadInterval time.Duration
	ticker         *time.Ticker
	reloadLock     sync.Mutex
}

func NewReloadableMap[K Hashable, T any](
	reloadInterval time.Duration,
	loader MapDataLoader[T],
	keyProvider KeyProvider[K, T],
) (*ReloadableMap[K, T], error) {
	m := &ReloadableMap[K, T]{
		loader:         loader,
		keyProvider:    keyProvider,
		reloadInterval: reloadInterval,
	}
	if loader != nil {
		err := m.loadData()
		if err != nil {
			return nil, err
		}

	}
	m.ticker = time.NewTicker(m.reloadInterval)
	go m.reloadDataTickerHandler()
	return m, nil
}

func (rm *ReloadableMap[K, T]) Store(element T) {
	rm.data.Store(rm.keyProvider(element), element)
}

func (rm *ReloadableMap[K, T]) Values() (values []T) {
	values = make([]T, 0)
	rm.data.Range(func(key, value any) bool {
		values = append(values, value.(T))
		return true
	})
	return
}

func (rm *ReloadableMap[K, T]) Get(key any) (*T, bool) {
	data, loaded := rm.data.Load(key)

	if !loaded {
		return nil, false
	}
	d := data.(T)
	return &d, loaded
}

func (rm *ReloadableMap[K, T]) Reload() {
	err := rm.loadData()
	if err != nil {
		return
	}
}

func (rm *ReloadableMap[K, T]) Range(hof func(key any, value T) bool) {
	rm.data.Range(func(k any, v any) bool {
		return hof(k, v.(T))
	})
}

func (rm *ReloadableMap[K, T]) reloadDataTickerHandler() {
	for {
		<-rm.ticker.C
		rm.ticker.Stop()
		_ = rm.loadData()
		rm.ticker.Reset(rm.reloadInterval)
	}
}

func (rm *ReloadableMap[K, T]) loadData() (err error) {
	rm.reloadLock.Lock()
	defer rm.reloadLock.Unlock()
	loadedData, err := rm.loader()
	if err != nil {
		return err
	}

	// Clean the map
	rm.data.Range(
		func(key, value interface{}) bool {
			rm.data.Delete(key)
			return true
		},
	)

	for _, m := range loadedData {
		rm.data.Store(rm.keyProvider(m), m)
	}

	return
}

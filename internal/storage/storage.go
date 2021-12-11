package storage

import (
	"sync"
)

type Storage struct {
	mx   sync.RWMutex
	keys map[string]interface{}
}

// New return instance of Storage
func New() *Storage {
	return &Storage{
		mx:   sync.RWMutex{},
		keys: map[string]interface{}{},
	}
}

func (s *Storage) Get(key string) (interface{}, bool) {
	s.mx.RLock()
	val, ok := s.keys[key]
	s.mx.RUnlock()

	return val, ok
}

func (s *Storage) Set(key string, value interface{}) error {
	s.mx.Lock()
	s.keys[key] = value
	s.mx.Unlock()

	return nil
}

func (s *Storage) Delete(key string) error {
	s.mx.Lock()
	delete(s.keys, key)
	s.mx.Unlock()

	return nil
}

func (s *Storage) Keys() []string {
	var i int

	res := make([]string, len(s.keys))
	for k := range s.keys {
		res[i] = k
		i++
	}

	return res
}

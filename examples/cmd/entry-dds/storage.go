package main

import (
	"errors"
	"sync"

	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/satori/go.uuid"
)

type entryStorage struct {
	mu   sync.RWMutex
	data map[string]microblag.Entry
}

func (s *entryStorage) Get(id string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	et, ok := s.data[id]
	if !ok {
		return nil, errors.New("Not found")
	}

	return &et, nil
}

func (s *entryStorage) Create(i interface{}) (string, error) {
	id := uuid.NewV1().String()

	entry, ok := i.(*microblag.Entry)
	if !ok {
		return "", errors.New("Unsupported type")
	}

	entry.ID = id

	s.mu.Lock()
	s.data[id] = *entry
	s.mu.Unlock()

	return id, nil
}

func (s *entryStorage) Delete(id string) error {

	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()

	return nil
}

func (s *entryStorage) Update(i interface{}) error {
	entry, ok := i.(*microblag.Entry)
	if !ok {
		return errors.New("Unsupported type")
	}

	s.mu.Lock()
	s.data[entry.ID] = *entry
	s.mu.Unlock()

	return nil
}

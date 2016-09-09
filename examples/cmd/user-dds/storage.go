package main

import (
	"errors"
	"log"
	"reflect"
	"sync"

	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/satori/go.uuid"
)

type userStorage struct {
	mu   sync.RWMutex
	data map[string]microblag.User
}

func (s *userStorage) Get(id string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	et, ok := s.data[id]
	if !ok {
		return nil, errors.New("Not found")
	}

	return &et, nil
}

func (s *userStorage) Create(i interface{}) (string, error) {
	id := uuid.NewV1().String()

	log.Printf("Creating: %s", id)

	user, ok := i.(*microblag.User)
	if !ok {
		return "", errors.New("Unsupported type: " + reflect.TypeOf(i).String())
	}

	user.ID = id

	s.mu.Lock()
	s.data[id] = *user
	s.mu.Unlock()

	return id, nil
}

func (s *userStorage) Delete(id string) error {
	log.Printf("Deleting: %s", id)

	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()

	return nil
}

func (s *userStorage) Update(i interface{}) error {
	user, ok := i.(*microblag.User)
	if !ok {
		return errors.New("Unsupported type")
	}

	s.mu.Lock()
	s.data[user.ID] = *user
	s.mu.Unlock()

	return nil
}

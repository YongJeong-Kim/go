package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"gogrpc/pb"
	"sync"
)

var ErrAlreadyExists = errors.New("already exists person")

type PersonStore interface {
	Save(person *pb.Person) error
}

type InMemoryPersonStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Person
}

func NewInMemoryPersonStore() *InMemoryPersonStore {
	return &InMemoryPersonStore{
		data: make(map[string]*pb.Person),
	}
}

func (store *InMemoryPersonStore) Save(person *pb.Person) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[person.Id] != nil {
		return ErrAlreadyExists
	}

	other := &pb.Person{}
	err := copier.Copy(other, person)
	if err != nil {
		return fmt.Errorf("cannot copy person. %w", err)
	}

	store.data[other.Id] = other
	return nil
}

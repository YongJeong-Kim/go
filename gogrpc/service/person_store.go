package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"gogrpc/pb"
	"log"
	"sync"
)

var ErrAlreadyExists = errors.New("already exists person")

type PersonStore interface {
	Save(person *pb.Person) error
	Find(id string) (*pb.Person, error)
	Search(filter *pb.Filter, found func(person *pb.Person) error) error
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

func (store *InMemoryPersonStore) Find(id string) (*pb.Person, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if p, ok := store.data[id]; ok {
		return p, nil
	} else {
		return nil, fmt.Errorf("Not Found Person\n")
	}
}

func (store *InMemoryPersonStore) Search(filter *pb.Filter, found func(person *pb.Person) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, person := range store.data {
		log.Print("filter brand:", filter.Brand)
		log.Print("person shirt brand:", person.Shirt.Brand)
		if filter.Brand == person.Shirt.Brand {
			other := &pb.Person{}
			err := copier.Copy(other, person)
			if err != nil {
				return nil
			}
			err = found(other)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

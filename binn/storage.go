package binn

import (
	"fmt"
	"time"
	"sync"

	"github.com/google/uuid"
)

const (
	MAX_CONTAINER_STORAGE_NUM_CONTAINER = 1000
	MAX_CONTAINER_STORAGE_NUM_ID = 1000
	MAX_EXPIRATION_HOUR = 10000
	MAX_MESSAGE_TEXT_LENGTH = 200
)

type ContainerKeeper interface {
	Get() (Container, error)
	Add(Container) error
}

type ContainerStorage struct {
	containers []Container
	idStorage  *IDStorage
	mux        *sync.Mutex
	validation bool
	expiration time.Duration
}

type IDStorage struct {
	ids map[string]time.Time
	mux *sync.Mutex
}

func NewContainerStorage(v bool, e time.Duration, s *IDStorage) *ContainerStorage {
	return &ContainerStorage{
		containers: []Container{},
		idStorage:  s,
		mux:        &sync.Mutex{},
		validation: v,
		expiration:	e,
	}
}

func DefaultContainerStorage() *ContainerStorage {
	return NewContainerStorage(true, 0, DefaultIDStorage())
}

func DefaultIDStorage() *IDStorage {
	return &IDStorage{
		ids: make(map[string]time.Time),
		mux: &sync.Mutex{},
	}
}

func (cs *ContainerStorage) Get() (Container, error) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	if len(cs.containers) == 0 {
		return nil, fmt.Errorf("this storage has no containers")
	}
	c := cs.containers[0]
	cs.containers = cs.containers[1:]
	if cs.expiration != 0 {
		d := time.Now().Add(cs.expiration)
		c = NewBottle(c.ID(), c.Message().Text, &d);
	} else {
		d := time.Now().Add(time.Duration(MAX_EXPIRATION_HOUR) * time.Hour)
		c = NewBottle(c.ID(), c.Message().Text, &d)
	}

	if cs.validation {
		cs.idStorage.Update(c.ID(), *c.ExpiredAt())
	}

	return c, nil
}

func (cs *ContainerStorage) Add(c Container) error {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	if cs.validation {
		if err := cs.idStorage.Use(c.ID()); err != nil {
			return err
		}
	}

	if len(cs.containers) >= MAX_CONTAINER_STORAGE_NUM_CONTAINER {
		cs.containers = cs.containers[1:]
	}

	messageText := c.Message().Text
	if len(c.Message().Text) > MAX_MESSAGE_TEXT_LENGTH {
		messageText = c.Message().Text[:MAX_MESSAGE_TEXT_LENGTH]
	}

	newID := GenerateID()
	if cs.validation {
		cs.idStorage.Add(newID, time.Now().Add(time.Duration(MAX_EXPIRATION_HOUR) * time.Hour))
	}

	c = NewBottle(newID, messageText, c.ExpiredAt())
	cs.containers = append(cs.containers, c)
	
	return nil
}

func (s *IDStorage) Add(id string, e time.Time) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.ids[id]; ok {
		return fmt.Errorf("this id (%#v) is already added", id)
	}
	s.ids[id] = e

	return nil
}

func (s *IDStorage) Use(id string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if v, ok := s.ids[id]; !ok {
		return fmt.Errorf("this id (%#v) is invalid", id)
	} else {
		if time.Now().After(v) {
			return fmt.Errorf("this id (%#v) is expired", id)
		}
	}

	delete(s.ids, id)
	
	return nil
}

func (s *IDStorage) Update(id string, e time.Time) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.ids[id]; !ok {
		return fmt.Errorf("this id (%#v) is not in storage", id)
	}
	s.ids[id] = e

	return nil
}

func GenerateID() string {
	uuidObj, _ := uuid.NewUUID()
	return uuidObj.String()
}

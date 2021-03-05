package main

import (
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/file"
)

type Storage struct {
	Store gokv.Store
	Dir   string
	Key   string
}

func CreateStore(dir, key string) *Storage {

	options := file.DefaultOptions
	options.Directory = dir

	// Create client
	client, err := file.NewStore(options)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	instance := &Storage{
		Store: client,
		Dir:   dir,
		Key:   key,
	}
	return instance

}

func (s *Storage) Save(item interface{}) error {
	return (s.Store.Set(s.Key, item))
}

func (s *Storage) Get(value interface{}) (bool, error) {
	found, err := s.Store.Get(s.Key, &value)
	return found, err
}

func (s *Storage) Delete() error {
	return (s.Store.Delete(s.Key))
}

func (s *Storage) Close() error {
	return (s.Store.Close())
}

package store

import (
	"fmt"
	"sync"
)

//store is a thread-safe in-memory store for the key-value pairs,
//multiple goroutines (one per client) will acces this
//this store concurrently, so every access to 'data' must go through
//the mutex - never access 'data' directly

type Store struct {
	mu   sync.RWMutex //create the mutex in the same struct
	data map[string]Value
}

// New creates an empty, ready-tp-use store
func New() *Store {
	return &Store{
		data: make(map[string]Value),
	}
}

// Set stores the given key-value pair in the store, overwritting any existing value for the key
// it recieves a pointer to the store, so it can modify the store's data in place
func (s *Store) Set(key, value string) {
	s.mu.Lock() //lock the mutex before accessing the data

	defer s.mu.Unlock() //unlock the mutex after accessing the data
	s.data[key] = StringValue{Value: value}
}

// Get returns the value for the key and whether it was found
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock() //shared lock for reading

	defer s.mu.RUnlock() // unlock
	value, ok := s.data[key]
	if !ok {
		return "", false
	}

	stringValue, ok := value.(StringValue)
	if !ok {
		return "", false
	}
	return stringValue.Value, true
}

// Delete removes the key-value pair from the store, if it exists
func (s *Store) Delete(key string) int {
	s.mu.Lock() //lock

	defer s.mu.Unlock()

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return 1
	}
	return 0
}

// Exists checks if the key exists in the store
func (s *Store) Exists(key string) bool {
	s.mu.RLock() //shared lock for reading

	defer s.mu.RUnlock() // unlock

	_, ok := s.data[key]
	return ok
}

// LPush pushes a value to the left (front) of a Redis List.
func (s *Store) LPush(key, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingValue, exists := s.data[key]

	// Case 1: the key does not exist yet.
	if !exists {
		list := ListValue{
			Values: []string{value},
		}

		s.data[key] = list
		return len(list.Values), nil
	}

	// Case 2 and 3: the key exists.
	list, ok := existingValue.(ListValue)
	if !ok {
		return 0, fmt.Errorf("WRONGTYPE operation against a key holding the wrong kind of value")
	}

	// Add the new value to the front of the list.
	list.Values = append([]string{value}, list.Values...)

	// Save the updated list back into the Store.
	s.data[key] = list

	return len(list.Values), nil
}

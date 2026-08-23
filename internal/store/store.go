package store

import (
	"sync"
)

//store is a thread-safe in-memory store for the key-value pairs,
//multiple goroutines (one per client) will acces this
//this store concurrently, so every access to 'data' must go through
//the mutex - never access 'data' directly

type store struct {
	mu   sync.RWMutex //create the mutex in the same struct
	data map[string]string
}

// New creates an empty, ready-tp-use store
func New() *store {
	return &store{
		data: make(map[string]string),
	}
}

// Set stores the given key-value pair in the store, overwritting any existing value for the key
// it recieves a pointer to the store, so it can modify the store's data in place
func (s *store) Set(key, value string) {
	s.mu.Lock() //lock the mutex before accessing the data

	defer s.mu.Unlock() //unlock the mutex after accessing the data
	s.data[key] = value
}

// Get returns the value for the key and whether it was found
func (s *store) Get(key string) (string, bool) {
	s.mu.RLock() //shared lock for reading

	defer s.mu.RUnlock() // unlock
	value, ok := s.data[key]
	return value, ok
}

// Delete removes the key-value pair from the store, if it exists
func (s *store) Delete(key string) int {
	s.mu.Lock() //lock

	defer s.mu.Unlock()

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return 1
	}
	return 0
}

// Exists checks if the key exists in the store
func (s *store) Exists(key string) bool {
	s.mu.RLock() //shared lock for reading

	defer s.mu.RUnlock() // unlock

	_, ok := s.data[key]
	return ok
}

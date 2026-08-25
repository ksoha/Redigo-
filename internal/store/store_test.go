package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestStoreSetGet(t *testing.T) {
	s := New()

	s.Set("foo", "bar")
	value, ok := s.Get("foo")

	if !ok {
		t.Errorf("Expected key 'foo' to exist")
	}

	if value != "bar" {
		t.Errorf("Expected value 'bar', got '%s'", value)
	}
}

// test for missing key
func TestStoreGetMissingKey(t *testing.T) {
	s := New()

	value, ok := s.Get("missing")
	if ok {
		t.Errorf("Expected key 'missing' to not exist")
	}
	if value != "" {
		t.Errorf("Expected value to be empty, got '%s'", value)
	}

}

// test for set overwriting existing key
func TestStoreSetOverwrite(t *testing.T) {
	s := New()

	s.Set("foo", "bar")
	s.Set("foo", "baz") //overwrites

	value, ok := s.Get("foo")
	if !ok {
		t.Errorf("Expected key 'foo' to exist")
	}
	if value != "baz" {
		t.Errorf("Expected value 'baz', got '%s'", value)
	}
}

// test for delete key
func TestStoreDelete(t *testing.T) {
	s := New()
	s.Set("foo", "bar")

	count := s.Delete("foo")
	if count != 1 {
		t.Errorf("Expected delete count to be 1, got %d", count)
	}

	_, ok := s.Get("foo")
	if ok {
		t.Errorf("Expected key 'foo' to not exist after deletion")
	}
}

// test for deleting the missing key
func TestStoreDeleteMissingKey(t *testing.T) {
	s := New()

	count := s.Delete("nonexistent")
	if count != 0 {
		t.Errorf("Expected Delete on missing key to return 0, got %d", count)
	}
}

// test if store exists
func TestStoreExists(t *testing.T) {
	s := New()
	s.Set("foo", "bar")

	if !s.Exists("foo") {
		t.Errorf("Expected key 'foo' to exist")
	}

	if s.Exists("missing") {
		t.Errorf("Expected key 'missing' to not exist")
	}
}

// test for concurrent access to the store
func TestStoreConcurrentAccess(t *testing.T) {
	s := New()

	var wg sync.WaitGroup

	//start multiple go routines
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d", i)

			s.Set(key, value)
			got, ok := s.Get(key)
			if !ok || got != value {
				t.Errorf("Expected key '%s' to have value '%s', got '%s'", key, value, got)
			}
		}(i)
	}
	wg.Wait()
}

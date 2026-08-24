package store

import "testing"

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

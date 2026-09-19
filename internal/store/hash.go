package store

import "fmt"

//HashValue represents a Redis Hash
//
//A Hash contains multiple key-value pairs, under one Redis Key
//
//ex :
//name -> "Soha"
//age -> 21
//city -> "Pune"

type HashValue struct {
	Fields map[string]string //Fields is a map that holds the key-value pairs of the hash
}

// Type returns the Redis data type
func (h HashValue) Type() string {
	return "hash"
}

// HSet sets a field inside a Redis Hash.
// It returns true if a new field was added.
func (s *Store) HSet(key string, field string, value string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingValue, exists := s.data[key]

	// If the key does not exist, create a new hash.
	if !exists {
		s.data[key] = HashValue{
			Fields: map[string]string{
				field: value,
			},
		}

		return true, nil
	}

	// The key exists, so check its datatype.
	hash, ok := existingValue.(HashValue)
	if !ok {
		return false, fmt.Errorf("WRONGTYPE operation against a key holding the wrong kind of value")
	}

	// Check whether this field already exists.
	_, fieldExists := hash.Fields[field]

	// Add or update the field.
	hash.Fields[field] = value

	// Store the updated hash back.
	s.data[key] = hash

	return !fieldExists, nil
}

// HGet retieves the value of a field in a redis hash
func (s *Store) HGet(key string, field string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingValue, exists := s.data[key]

	//if the redis ket doesnt not exist
	if !exists {
		return "", false, nil
	}

	// The key exists, so check whether it is a HashValue.
	hash, ok := existingValue.(HashValue)
	if !ok {
		return "", false, fmt.Errorf("WRONGTYPE operation against a key holding the wrong kind of value")
	}

	// Look for the requested field.
	value, fieldExists := hash.Fields[field]

	if !fieldExists {
		return "", false, nil
	}

	return value, true, nil
}

package store

// SetValue represents a Redis Set.
//
// A Set contains unique members under one Redis key.
//
// Example:
//
// "Go"
// "Redis"
// "Docker"
type SetValue struct {
	Members map[string]struct{}
}

// Type returns the Redis data type.
func (s SetValue) Type() string {
	return "set"
}

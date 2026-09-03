package store

// ListValue represents redis list
type ListValue struct {
	Values []string
}

// Type returns the Redis data type
func (l ListValue) Type() string {
	return "list"
}

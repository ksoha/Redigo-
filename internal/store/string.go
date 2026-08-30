package store

// StringValue represensts a redis string value, which is a simple string
type StringValue struct {
	Value string
}

// Type returns the redis data type
func (s StringValue) Type() string {
	return "string"
}

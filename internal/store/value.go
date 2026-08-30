package store

// Value represents any value that can be stored in redigo
type Value interface {
	Type() string
}

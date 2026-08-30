package store

type StringValue struct {
	Value string
}

func (s StringValue) Type() string {
	return "string"
}

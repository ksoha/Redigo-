package store

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

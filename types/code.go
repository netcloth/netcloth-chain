package types

import "fmt"

// Code & Storage
type (
	Code []byte

	// Storage is account storage type alias
	Storage map[Hash]Hash
)

func (c Code) String() string {
	return string(c)
}

func (s Storage) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}
	return
}

// Copy returns a copy of storage
func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}

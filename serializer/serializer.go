package serializer

import (
	"encoding/json"
	"log"
)

type Serializer interface {
	Serialize(object interface{}) []byte
}

type SerializerFunc func(interface{}) []byte

// Serialize implements the Handler interface
func (s SerializerFunc) Serialize(o interface{}) []byte {
	return s(o)
}

func Serialize(object interface{}) []byte {
	o, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("error marshaling %#v: %s", object, err)
	}
	return o
}

func NewSerializer() Serializer {
	return SerializerFunc(Serialize)
}

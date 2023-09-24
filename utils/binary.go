package utils

import (
	"bytes"
	"encoding/gob"
)

func EncodeToBinary(obj any) (encoded []byte, err error) {
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)

	err = enc.Encode(obj)
	if err != nil {
		return encoded, err
	}
	return buffer.Bytes(), err
}

func DecodeFromBinary[T any](data []byte) (obj *T, err error) {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	obj = new(T)
	err = dec.Decode(obj)
	return
}

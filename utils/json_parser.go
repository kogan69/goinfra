package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	UseNumber bool = true
)

func FromJSON[T any](r io.Reader, useNumber bool) (obj *T, err error) {
	if r == nil {
		return
	}

	decoder := json.NewDecoder(r)
	if useNumber {
		decoder.UseNumber()
	}
	obj = new(T)
	err = decoder.Decode(obj)

	return
}

func FromJsonData[T any](data []byte, useNumber bool) (obj *T, err error) {
	bb := bytes.NewBuffer(data)
	return FromJSON[T](bb, useNumber)
}

func ToJSON[T any](obj T) (data string, err error) {
	bb := new(bytes.Buffer)
	enc := json.NewEncoder(bb)
	err = enc.Encode(obj)
	if err != nil {
		return "", err
	}
	data = bb.String()
	return
}

func ToJsonAsByteArray[T any](obj T) (data []byte, err error) {
	bb := new(bytes.Buffer)
	enc := json.NewEncoder(bb)
	err = enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	data = bb.Bytes()
	return
}

func ToJsonAsByteBuffer[T any](obj T) (bb *bytes.Buffer, err error) {
	bb = new(bytes.Buffer)
	enc := json.NewEncoder(bb)
	err = enc.Encode(obj)
	return
}

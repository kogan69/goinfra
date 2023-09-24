package utils

import (
	"github.com/gofrs/uuid/v5"
)

func StringRef(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StringVal(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func Ref[T any](val T) *T {
	return &val
}

func UuidOrNil(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

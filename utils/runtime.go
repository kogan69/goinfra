package utils

import (
	"runtime"
)

func FunctionName() string {
	pcs := make([]uintptr, 2)
	depth := runtime.Callers(1, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	f, again := frames.Next()
	if !again {
		return ""
	}
	f, _ = frames.Next()
	return f.Function
}

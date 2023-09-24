package utils

import (
	"net/http"
)

func NewHttpClient(followRedirect bool) *http.Client {
	client := http.DefaultClient
	if followRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return client
}

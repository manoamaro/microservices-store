package infra

import (
	"net/http"
)

type HttpService struct {
	Host   string
	Client *http.Client
}

func NewHttpService(host string) *HttpService {
	return &HttpService{
		Host:   host,
		Client: &http.Client{},
	}
}

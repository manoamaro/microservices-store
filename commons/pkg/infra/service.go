package infra

import (
	"net/http"
)

type Service struct {
	Host   string
	Client *http.Client
}

func NewService(host string) *Service {
	return &Service{
		Host:   host,
		Client: &http.Client{},
	}
}

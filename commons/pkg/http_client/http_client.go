package http_client

import (
	"net/http"
	"time"
)

type HttpClient struct {
	Host   string
	Client *http.Client
}

func NewHttpClient(host string) *HttpClient {
	return &HttpClient{
		Host: host,
		Client: &http.Client{
			Transport: &http.Transport{},
			Timeout:   time.Second * 60,
		},
	}
}

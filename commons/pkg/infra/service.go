package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/manoamaro/gobreaker/v2"
)

type Endpoint[Res any] struct {
	Method  string
	service *Service
	CB      *gobreaker.CircuitBreaker[Res]
}

func NewEndpoint[T any](service *Service, method string, maxRequests uint32, interval time.Duration) *Endpoint[T] {
	st := gobreaker.Settings{
		Name:        "test",
		MaxRequests: maxRequests,
		Interval:    interval,
	}

	return &Endpoint[T]{
		Method:  method,
		service: service,
		CB:      gobreaker.NewCircuitBreaker[T](st),
	}
}

func (e *Endpoint[Res]) Execute(path string, headers map[string]string, body any) (Res, error) {
	var reqBody []byte
	var response Res

	if body != nil {
		if _reqBody, err := json.Marshal(body); err != nil {
			return response, err
		} else {
			reqBody = _reqBody
		}
	}

	fullPath := fmt.Sprintf("%s%s", e.service.Host, path)

	request, err := http.NewRequest(e.Method, fullPath, bytes.NewReader(reqBody))
	if err != nil {
		return response, err
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	response, err = e.CB.Execute(func() (Res, error) {
		var r Res
		if response, err := e.service.Client.Do(request); err != nil {
			return r, err
		} else if response.StatusCode != http.StatusOK {
			return r, fmt.Errorf("error fetching inventory")
		} else {
			defer response.Body.Close()
			if body, err := io.ReadAll(response.Body); err != nil {
				return r, err
			} else {
				var res Res
				if err := json.Unmarshal(body, &res); err != nil {
					return r, err
				}
				return res, nil
			}
		}
	})
	if err != nil {
		return response, err
	} else {
		return response, nil
	}
}

type IService interface {
}

type Service struct {
	IService
	Host   string
	Client *http.Client
}

func NewService(host string) *Service {
	return &Service{
		Host:   host,
		Client: &http.Client{},
	}
}

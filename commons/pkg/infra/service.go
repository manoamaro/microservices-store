package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/manoamaro/gobreaker/v2"
)

type Endpoint[Res any] struct {
	method  string
	path    string
	service *Service
	CB      *gobreaker.CircuitBreaker[Res]
}

type EndpointExec[Res any] struct {
	*Endpoint[Res]
	headers     map[string]string
	queryParams map[string]string
	pathParams  map[string]string
	body        any
}

func NewEndpoint[T any](service *Service, method string, path string, maxRequests uint32, interval time.Duration) *Endpoint[T] {
	st := gobreaker.Settings{
		Name:        "test",
		MaxRequests: maxRequests,
		Interval:    interval,
	}

	return &Endpoint[T]{
		method:  method,
		path:    path,
		service: service,
		CB:      gobreaker.NewCircuitBreaker[T](st),
	}
}

func (e *Endpoint[Res]) Start() *EndpointExec[Res] {
	return &EndpointExec[Res]{
		Endpoint:    e,
		headers:     map[string]string{},
		queryParams: map[string]string{},
		pathParams:  map[string]string{},
	}
}

func (e *EndpointExec[Res]) WithHeaders(headers map[string]string) *EndpointExec[Res] {
	for k, v := range headers {
		e.headers[k] = v
	}
	return e
}

func (e *EndpointExec[Res]) WithAuthorization(token string) *EndpointExec[Res] {
	e.WithHeaders(map[string]string{"Authorization": token})
	return e
}

func (e *EndpointExec[Res]) WithQueryParams(queryParams map[string]string) *EndpointExec[Res] {
	e.queryParams = queryParams
	return e
}

func (e *EndpointExec[Res]) WithPathParams(pathParams map[string]string) *EndpointExec[Res] {
	e.pathParams = pathParams
	return e
}

func (e *EndpointExec[Res]) WithPathParam(name, value string) *EndpointExec[Res] {
	e.pathParams[name] = value
	return e
}

func (e *EndpointExec[Res]) WithBody(body any) *EndpointExec[Res] {
	e.body = body
	return e
}

func (e *EndpointExec[Res]) Execute() (Res, error) {
	var reqBody []byte
	var response Res

	if e.body != nil {
		if _reqBody, err := json.Marshal(e.body); err != nil {
			return response, err
		} else {
			reqBody = _reqBody
		}
	}

	fullPath := fmt.Sprintf("%s%s", e.service.Host, e.path)

	for k, v := range e.pathParams {
		fullPath = strings.ReplaceAll(fullPath, k, v)
	}

	_url, err := url.Parse(fullPath)
	if err != nil {
		return response, err
	}
	for k, v := range e.queryParams {
		_url.Query().Add(k, v)
	}

	req, err := http.NewRequest(e.method, _url.String(), bytes.NewReader(reqBody))
	if err != nil {
		return response, err
	}

	for k, v := range e.headers {
		req.Header.Add(k, v)
	}

	response, err = e.CB.Execute(func() (Res, error) {
		var r Res
		if response, err := e.service.Client.Do(req); err != nil {
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

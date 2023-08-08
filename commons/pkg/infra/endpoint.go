package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/infra/cb"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Endpoint[Res any] struct {
	method  string
	path    string
	service *HttpService
	CB      *cb.CircuitBreaker[Res]
}

func NewEndpoint[T any](
	service *HttpService,
	method string,
	path string,
	maxRequests int,
	interval time.Duration,
) *Endpoint[T] {
	return &Endpoint[T]{
		method:  method,
		path:    path,
		service: service,
		CB:      cb.NewCircuitBreaker[T](maxRequests, interval),
	}
}

type RequestEndpointCommand[Res any] struct {
	*Endpoint[Res]
	headers     map[string]string
	queryParams map[string]string
	pathParams  map[string]string
	body        any
}

func (e *Endpoint[Res]) Start() *RequestEndpointCommand[Res] {
	return &RequestEndpointCommand[Res]{
		Endpoint:    e,
		headers:     map[string]string{},
		queryParams: map[string]string{},
		pathParams:  map[string]string{},
	}
}

func (e *RequestEndpointCommand[Res]) WithHeader(key, value string) *RequestEndpointCommand[Res] {
	e.headers[key] = value
	return e
}

func (e *RequestEndpointCommand[Res]) WithAuthorization(token string) *RequestEndpointCommand[Res] {
	e.WithHeader("Authorization", token)
	return e
}

func (e *RequestEndpointCommand[Res]) WithQueryParam(key, value string) *RequestEndpointCommand[Res] {
	e.queryParams[key] = value
	return e
}

func (e *RequestEndpointCommand[Res]) WithPathParam(name, value string) *RequestEndpointCommand[Res] {
	e.pathParams[name] = value
	return e
}

func (e *RequestEndpointCommand[Res]) WithBody(body any) *RequestEndpointCommand[Res] {
	e.body = body
	return e
}

func (e *RequestEndpointCommand[Res]) Execute() (Res, error) {
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

	response, err = e.CB.Call(func() (Res, error) {
		var r Res
		if response, err := e.service.Client.Do(req); err != nil {
			return r, err
		} else if response.StatusCode != http.StatusOK {
			return r, fmt.Errorf(response.Status)
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

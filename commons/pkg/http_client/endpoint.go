package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/cb"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Endpoint[Req any, Res any] struct {
	method  string
	path    string
	service *HttpClient
	CB      *cb.CircuitBreaker[Res]
}

func NewEndpoint[Req any, Res any](
	service *HttpClient,
	method string,
	path string,
	maxRequests int,
	interval time.Duration,
) *Endpoint[Req, Res] {
	return &Endpoint[Req, Res]{
		method:  method,
		path:    path,
		service: service,
		CB:      cb.NewCircuitBreaker[Res](maxRequests, interval),
	}
}

type RequestEndpointCommand[Req, Res any] struct {
	*Endpoint[Req, Res]
	headers     map[string]string
	queryParams map[string]string
	pathParams  map[string]string
	body        Req
	hasBody     bool
}

func (e *Endpoint[Req, Res]) Start() *RequestEndpointCommand[Req, Res] {
	return &RequestEndpointCommand[Req, Res]{
		Endpoint:    e,
		headers:     map[string]string{},
		queryParams: map[string]string{},
		pathParams:  map[string]string{},
		hasBody:     false,
	}
}

func (e *RequestEndpointCommand[Req, Res]) WithHeader(key, value string) *RequestEndpointCommand[Req, Res] {
	e.headers[key] = value
	return e
}

func (e *RequestEndpointCommand[Req, Res]) WithAuthorization(token string) *RequestEndpointCommand[Req, Res] {
	e.WithHeader("Authorization", token)
	return e
}

func (e *RequestEndpointCommand[Req, Res]) WithQueryParam(key, value string) *RequestEndpointCommand[Req, Res] {
	e.queryParams[key] = value
	return e
}

func (e *RequestEndpointCommand[Req, Res]) WithPathParam(name, value string) *RequestEndpointCommand[Req, Res] {
	e.pathParams[name] = value
	return e
}

func (e *RequestEndpointCommand[Req, Res]) WithBody(body Req) *RequestEndpointCommand[Req, Res] {
	e.body = body
	e.hasBody = true
	return e
}

func (e *RequestEndpointCommand[Req, Res]) Execute() (Res, error) {
	var reqBody []byte
	var response Res

	if e.hasBody {
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

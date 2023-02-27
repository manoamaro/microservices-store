package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type IService interface {
}

type Service struct {
	Host   string
	Client *http.Client
	CB     *gobreaker.CircuitBreaker
}

func NewService(host string, name string, maxRequests int, interval int) Service {
	var st gobreaker.Settings
	st.Name = name
	st.MaxRequests = uint32(maxRequests)
	st.Interval = time.Duration(interval)

	return Service{
		Host:   host,
		CB:     gobreaker.NewCircuitBreaker(st),
		Client: &http.Client{},
	}
}

func Req[T any](client *http.Client, method string, path string, body any) (*T, error) {

	var reqBody []byte
	if body != nil {
		if _reqBody, err := json.Marshal(body); err != nil {
			return nil, err
		} else {
			reqBody = _reqBody
		}
	}

	if request, err := http.NewRequest(method, path, bytes.NewReader(reqBody)); err != nil {
		return nil, err
	} else if response, err := client.Do(request); err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching inventory")
	} else {
		defer response.Body.Close()
		if body, err := ioutil.ReadAll(response.Body); err != nil {
			return nil, err
		} else {
			var res T
			if err := json.Unmarshal(body, &res); err != nil {
				return nil, err
			}
			return &res, nil
		}
	}
}

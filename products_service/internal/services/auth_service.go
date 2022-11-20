package services

import (
	"fmt"
	"github.com/sony/gobreaker"
	"net/http"
)

type AuthService interface {
	Validate(token string, audiences []string) (error, bool)
}

type DefaultAuthService struct {
	host   string
	client *http.Client
	cb     *gobreaker.CircuitBreaker
}

func NewDefaultAuthService(host string) AuthService {
	var st gobreaker.Settings
	st.Name = "AuthService"
	st.MaxRequests = 10
	st.Interval = 3000

	return &DefaultAuthService{
		host:   host,
		cb:     gobreaker.NewCircuitBreaker(st),
		client: &http.Client{},
	}
}

func (d *DefaultAuthService) Validate(token string, audiences []string) (error, bool) {
	response, err := d.cb.Execute(func() (interface{}, error) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/auth/verify", d.host), nil)
		if err != nil {
			return false, err
		}

		response, err := d.client.Do(request)
		if err != nil {
			return false, err
		}

		if response.StatusCode != http.StatusOK {
			return false, nil
		}

		return true, nil
	})

	return err, response.(bool)
}

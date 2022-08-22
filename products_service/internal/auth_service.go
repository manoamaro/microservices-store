package internal

import (
	"errors"
	"net/http"
)

type AuthService interface {
	Validate(token string) (error, bool)
}

func NewAuthService() AuthService {
	host := GetEnv("AUTH_URL", "http://localhost:8081/auth")
	return &DefaultAuthService{
		host:   host,
		client: &http.Client{},
	}
}

type DefaultAuthService struct {
	host   string
	client *http.Client
}

func (d *DefaultAuthService) Validate(token string) (error, bool) {
	request, _ := http.NewRequest(http.MethodGet, d.host+"/validate", nil)
	request.Header.Add("Authorization", token)
	response, err := d.client.Do(request)
	if err != nil {
		return err, false
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status), false
	}
	return nil, true
}

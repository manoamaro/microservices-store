package services

import (
	"encoding/json"
	"errors"
	"manoamaro.github.com/commons/pkg/collections"
	"net/http"
)

type AuthService interface {
	Validate(token string, audiences []string) (error, bool)
}

type DefaultAuthService struct {
	host   string
	client *http.Client
}

func NewDefaultAuthService(host string) AuthService {
	return &DefaultAuthService{
		host:   host,
		client: &http.Client{},
	}
}

func (d *DefaultAuthService) Validate(token string, audiences []string) (error, bool) {
	request, _ := http.NewRequest(http.MethodGet, d.host+"/validate", nil)
	request.Header.Add("Authorization", token)
	response, err := d.client.Do(request)
	if err != nil {
		return err, false
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status), false
	}
	body := struct {
		Audiences []string `json:"audiences"`
		Flags     []string `json:"flags"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return err, false
	}
	return nil, len(audiences) == 0 || collections.ContainsAny(body.Audiences, audiences)
}

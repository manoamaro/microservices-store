package services

import (
	"fmt"
	"net/http"

	"github.com/manoamaro/microservices-store/commons/pkg/infra"
)

type AuthService interface {
	infra.IService
	Validate(token string, audiences []string) (error, bool)
}

type DefaultAuthService struct {
	infra.Service
}

func NewDefaultAuthService(host string) AuthService {
	return &DefaultAuthService{
		infra.NewService(host, "AuthService", 10, 3000),
	}
}

func (d *DefaultAuthService) Validate(token string, audiences []string) (error, bool) {
	response, err := d.CB.Execute(func() (interface{}, error) {

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/auth/verify", d.Host), nil)
		if err != nil {
			return false, err
		}

		request.Header.Add("Authorization", token)

		response, err := d.Client.Do(request)
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

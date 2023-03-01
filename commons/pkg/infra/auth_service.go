package infra

import (
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"net/http"
)

type AuthService interface {
	IService
	Validate(token string, audiences ...string) (bool, error)
}

type VerifyResponse struct {
	Audiences []string `json:"audiences"`
	Flags     []string `json:"flags"`
}

type DefaultAuthService struct {
	*Service
	verifyEndpoint *Endpoint[VerifyResponse]
}

func NewDefaultAuthService(host string) AuthService {
	service := NewService(host)
	return &DefaultAuthService{
		Service:        service,
		verifyEndpoint: NewEndpoint[VerifyResponse](service, http.MethodGet, 10, 10*10e9),
	}
}

func (d *DefaultAuthService) Validate(token string, audiences ...string) (bool, error) {
	response, err := d.verifyEndpoint.Execute("/public/verify", map[string]string{"Authorization": token}, nil)
	if err != nil {
		return false, err
	}
	return collections.ContainsAny(audiences, response.Audiences), err
}

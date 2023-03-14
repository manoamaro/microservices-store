package infra

import (
	"net/http"
)

type AuthService interface {
	Validate(token string, audiences ...string) (*VerifyResponse, error)
}

type VerifyResponse struct {
	UserId    string   `json:"user_id"`
	Audiences []string `json:"audiences"`
	Flags     []string `json:"flags"`
}

type httpAuthService struct {
	*Service
	verifyEndpoint *Endpoint[VerifyResponse]
}

func NewDefaultAuthService(host string) AuthService {
	service := NewService(host)
	return &httpAuthService{
		Service:        service,
		verifyEndpoint: NewEndpoint[VerifyResponse](service, http.MethodGet, "/public/verify", 10, 10*10e9),
	}
}

func (d *httpAuthService) Validate(token string, audiences ...string) (*VerifyResponse, error) {
	response, err := d.verifyEndpoint.Start().
		WithAuthorization(token).
		Execute()
	if err != nil {
		return nil, err
	}
	return &response, err
}

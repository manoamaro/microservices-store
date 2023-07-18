package infra

import (
	"errors"
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"net/http"
	"time"
)

var (
	ErrorUnauthorized = errors.New("unauthorized")
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
	verifyEndpoint *Endpoint[VerifyResponse]
}

func NewHttpAuthService(host string) AuthService {
	service := NewHttpService(host)
	return &httpAuthService{
		verifyEndpoint: NewEndpoint[VerifyResponse](
			service,
			http.MethodGet,
			"/public/verify",
			10,
			time.Second*60,
		),
	}
}

func (d *httpAuthService) Validate(token string, audiences ...string) (*VerifyResponse, error) {
	response, err := d.verifyEndpoint.Start().
		WithAuthorization(token).
		Execute()
	if err != nil {
		return nil, err
	}
	if response.Audiences == nil || len(response.Audiences) == 0 {
		return &response, nil
	} else if collections.ContainsAny(response.Audiences, audiences) {
		return &response, nil
	} else {
		return nil, ErrorUnauthorized
	}
}

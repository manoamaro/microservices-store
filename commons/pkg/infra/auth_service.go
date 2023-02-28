package infra

import (
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"net/http"
)

type AuthService interface {
	IService
	Validate(token string, audiences ...string) (error, bool)
}

type DefaultAuthService struct {
	Service
}

func NewDefaultAuthService(host string) AuthService {
	return &DefaultAuthService{
		NewService(host, "AuthService", 10, 3000),
	}
}

type VerifyResponse struct {
	Audiences []string `json:"audiences"`
	Flags     []string `json:"flags"`
}

func (d *DefaultAuthService) Validate(token string, audiences ...string) (error, bool) {
	response, err := d.CB.Execute(func() (interface{}, error) {

		res, err := Req[VerifyResponse](d.Client, http.MethodGet, fmt.Sprintf("%s/public/verify", d.Host), map[string]string{"Authorization": token}, nil)
		if err != nil {
			return nil, err
		}

		return collections.ContainsAny(audiences, res.Audiences), nil
	})

	return err, response.(bool)
}

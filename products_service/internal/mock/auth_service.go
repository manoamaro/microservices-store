package mock

import (
	"github.com/manoamaro/microservices-store/commons/pkg/infra"
)

type AuthService struct {
}

func (m *AuthService) Validate(token string, audiences ...string) (*infra.VerifyResponse, error) {
	return nil, nil
}

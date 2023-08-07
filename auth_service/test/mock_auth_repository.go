package test

import (
	"github.com/manoamaro/microservices-store/auth_service/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{}
}

func (m *MockAuthRepository) Get(id uint) (auth models.Auth, err error) {
	args := m.Called(id)
	return args.Get(0).(models.Auth), args.Error(1)
}

func (m *MockAuthRepository) Create(email string, plainPassword string, audience []string, flags []string) (auth *models.Auth, err error) {
	args := m.Called(email, plainPassword, audience, flags)
	return args.Get(0).(*models.Auth), args.Error(1)
}

func (m *MockAuthRepository) Authenticate(email string, plainPassword string) (auth *models.Auth, found bool) {
	args := m.Called(email, plainPassword)
	return args.Get(0).(*models.Auth), args.Bool(1)
}

func (m *MockAuthRepository) InvalidateToken(token *models.UserClaims, rawToken string) error {
	args := m.Called(token, rawToken)
	return args.Error(0)
}

func (m *MockAuthRepository) IsInvalidatedToken(rawToken string) bool {
	args := m.Called(rawToken)
	return args.Bool(0)
}

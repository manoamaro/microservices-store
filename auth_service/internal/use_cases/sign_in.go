package use_cases

import (
	"fmt"
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
)

type SignInDTO struct {
	Email         string
	PlainPassword string
}

type SignInResultDTO struct {
	Token        string
	RefreshToken string
}

type SignInUseCase interface {
	SignIn(SignInDTO) (SignInResultDTO, error)
}

type signInUseCase struct {
	repository repositories.AuthRepository
}

func NewSignInUseCase(repository repositories.AuthRepository) SignInUseCase {
	return &signInUseCase{
		repository: repository,
	}
}

func (s *signInUseCase) SignIn(signInDTO SignInDTO) (SignInResultDTO, error) {
	result := SignInResultDTO{}
	if auth, found := s.repository.Authenticate(signInDTO.Email, signInDTO.PlainPassword); !found {
		return result, fmt.Errorf("user not found")
	} else if accessToken, refreshToken, err := helpers.CreateTokens(auth.ID, auth.DomainArray(), auth.FlagsArray()); err != nil {
		return result, err
	} else {
		result.Token = accessToken
		result.RefreshToken = refreshToken
		return result, nil
	}
}

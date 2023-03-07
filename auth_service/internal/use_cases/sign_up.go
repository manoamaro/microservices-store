package use_cases

import (
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
)

type SignUpDTO struct {
	Email         string
	PlainPassword string
	Audience      []string
	Flags         []string
}

type SignUpResultDTO struct {
	Token        string
	RefreshToken string
}

type SignUpUseCase interface {
	SignUp(SignUpDTO) (SignUpResultDTO, error)
}

type signUpUseCase struct {
	repository repositories.AuthRepository
}

func NewSignUpUseCase(repository repositories.AuthRepository) SignUpUseCase {
	return &signUpUseCase{repository: repository}
}

func (s *signUpUseCase) SignUp(signUpDTO SignUpDTO) (SignUpResultDTO, error) {
	result := SignUpResultDTO{}
	if auth, err := s.repository.Create(signUpDTO.Email, signUpDTO.PlainPassword, signUpDTO.Audience, signUpDTO.Flags); err != nil {
		return result, err
	} else if accessToken, refreshToken, err := helpers.CreateTokens(auth.ID, auth.DomainArray(), auth.FlagsArray()); err != nil {
		return result, err
	} else {
		result.Token = accessToken
		result.RefreshToken = refreshToken
		return result, nil
	}
}

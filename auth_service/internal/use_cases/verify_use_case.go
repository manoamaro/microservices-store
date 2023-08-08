package use_cases

import (
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
)

type VerifyDTO struct {
	Token string
}

type VerifyResultDTO struct {
	ID       string
	Audience []string
	Flags    []string
}

type VerifyUseCase interface {
	Verify(VerifyDTO) (VerifyResultDTO, error)
}

type verifyUseCase struct {
	repository repositories.AuthRepository
}

func NewVerifyUseCase(repository repositories.AuthRepository) VerifyUseCase {
	return &verifyUseCase{repository: repository}
}

func (v *verifyUseCase) Verify(dto VerifyDTO) (VerifyResultDTO, error) {
	var result VerifyResultDTO
	if userClaims, err := helpers.GetClaimsFromToken(dto.Token); err != nil {
		return result, err
	} else if v.repository.IsInvalidatedToken(dto.Token) {
		return result, ErrTokenInvalidated
	} else {
		result.ID = userClaims.ID
		result.Audience = userClaims.Audience
		result.Flags = userClaims.Flags
		return result, nil
	}
}

package use_cases

import (
	"errors"
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/repositories"
	"strconv"
)

var ErrTokenInvalidated = errors.New("token invalidated")

type RefreshTokenDTO struct {
	RefreshToken string
}

type RefreshTokenResultDTO struct {
	Token        string
	RefreshToken string
}

type RefreshTokenUseCase interface {
	RefreshToken(RefreshTokenDTO) (RefreshTokenResultDTO, error)
}

type refreshTokenUseCase struct {
	repository repositories.AuthRepository
}

func NewRefreshTokenUseCase(repository repositories.AuthRepository) RefreshTokenUseCase {
	return &refreshTokenUseCase{repository: repository}
}

func (r *refreshTokenUseCase) RefreshToken(dto RefreshTokenDTO) (RefreshTokenResultDTO, error) {
	var result RefreshTokenResultDTO
	if claims, err := helpers.GetClaimsFromRefreshToken(dto.RefreshToken); err != nil {
		return result, err
	} else if r.repository.IsInvalidatedToken(dto.RefreshToken) {
		return result, ErrTokenInvalidated
	} else if authId, err := strconv.ParseUint(claims.ID, 10, 32); err != nil {
		return result, err
	} else if auth, err := r.repository.Get(uint(authId)); err != nil {
		return result, err
	} else if token, refreshToken, err := helpers.CreateTokens(auth.ID, auth.DomainArray(), auth.FlagsArray()); err != nil {
		return result, err
	} else {
		result.Token = token
		result.RefreshToken = refreshToken
		return result, err
	}
}

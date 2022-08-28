package grpc_service

import (
	"context"
	"errors"
	"manoamaro.github.com/auth_service/internal/repositories"
	"manoamaro.github.com/commons/pkg/services"
)

type AuthServiceServer struct {
	repository repositories.AuthRepository
}

func (a *AuthServiceServer) Verify(ctx context.Context, request *services.VerifyRequest) (*services.AuthResponse, error) {
	if !a.repository.CheckToken(request.RawToken) {
		return nil, errors.New("token not valid")
	}

	claims, err := a.repository.GetClaimsFromToken(request.RawToken)
	if err != nil {
		return nil, err
	}

	return &services.AuthResponse{
		Audiences: claims.Audience,
		Flags:     claims.Flags,
	}, nil
}

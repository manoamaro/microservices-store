package use_cases

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/use_cases"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var signInUseCase = use_cases.NewSignInUseCase(mockAuthRepository)

func Test_signInUseCase_SignInWithValidCredentials(t *testing.T) {
	mockAuthRepository.On("Authenticate", validAuth.Email, validAuth.Password).Return(&validAuth, true).Once()
	args := use_cases.SignInDTO{
		Email:         validAuth.Email,
		PlainPassword: validAuth.Password,
	}
	var domains jwt.ClaimStrings = validAuth.DomainArray()

	resultDTO, err := signInUseCase.SignIn(args)
	assert.NoError(t, err)
	assert.NotEmpty(t, resultDTO.Token)
	assert.NotEmpty(t, resultDTO.RefreshToken)
	userClaims, err := helpers.GetClaimsFromToken(resultDTO.Token)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(int(validAuth.ID)), userClaims.ID)
	assert.Equal(t, domains, userClaims.Audience)
	assert.Equal(t, validAuth.FlagsArray(), userClaims.Flags)
	mockAuthRepository.AssertExpectations(t)
}

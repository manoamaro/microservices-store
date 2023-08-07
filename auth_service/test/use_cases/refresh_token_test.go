package use_cases

import (
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/use_cases"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"github.com/manoamaro/microservices-store/auth_service/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

var mockAuthRepository = test.NewMockAuthRepository()
var useCase = use_cases.NewRefreshTokenUseCase(mockAuthRepository)

var flag1 = models.Flag{
	Model: gorm.Model{ID: 1},
	Name:  "flag1",
}
var domain1 = models.Domain{
	Model:  gorm.Model{ID: 1},
	Domain: "domain1",
}

var validAuth = models.Auth{
	Model:    gorm.Model{ID: 1},
	Email:    "validAuth@example.com",
	Password: "validAuthPassword",
	Salt:     "1",
	Flags:    []models.Flag{flag1},
	Domains:  []models.Domain{domain1},
}

func Test_refreshTokenUseCase_RefreshTokenWhenTokenIsValid(t *testing.T) {
	// Setup
	token, refreshToken, _ := helpers.CreateTokens(validAuth.ID, validAuth.DomainArray(), validAuth.FlagsArray())
	mockAuthRepository.On("Get", validAuth.ID).Return(validAuth, nil)
	mockAuthRepository.On("IsInvalidatedToken", token).Return(false)
	mockAuthRepository.On("IsInvalidatedToken", refreshToken).Return(false)
	arg := use_cases.RefreshTokenDTO{
		RefreshToken: refreshToken,
	}

	// Run the test
	// Need to sleep so the generated token is different
	time.Sleep(1 * time.Second)
	result, err := useCase.RefreshToken(arg)

	// Verify the results
	assert.NoError(t, err)
	assert.NotEqualf(t, token, result.Token, "new token should be different")
	assert.NotEqualf(t, refreshToken, result.RefreshToken, "new refresh token should be different")
	newUserClaims, err := helpers.GetClaimsFromToken(result.Token)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(int(validAuth.ID)), newUserClaims.ID)
}

func Test_refreshTokenUseCase_RefreshTokenWhenTokenIsInvalid(t *testing.T) {
	// Setup
	_, refreshToken, _ := helpers.CreateTokens(validAuth.ID, validAuth.DomainArray(), validAuth.FlagsArray())
	mockAuthRepository.On("Get", validAuth.ID).Return(validAuth, nil)
	mockAuthRepository.On("IsInvalidatedToken", refreshToken).Return(true)
	arg := use_cases.RefreshTokenDTO{
		RefreshToken: refreshToken,
	}

	// Run the test
	_, err := useCase.RefreshToken(arg)

	// Verify the results
	assert.Error(t, err, use_cases.ErrTokenInvalidated)
}

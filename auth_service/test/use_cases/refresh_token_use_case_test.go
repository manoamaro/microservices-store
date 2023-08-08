package use_cases

import (
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/use_cases"
	"github.com/manoamaro/microservices-store/auth_service/models"
	"github.com/manoamaro/microservices-store/auth_service/test"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

type RefreshTokenUseCaseTestSuite struct {
	suite.Suite
	authRepository *test.MockAuthRepository
	useCase        use_cases.RefreshTokenUseCase
}

func TestRefreshTokenUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenUseCaseTestSuite))
}

func (suite *RefreshTokenUseCaseTestSuite) SetupTest() {
	suite.authRepository = test.NewMockAuthRepository()
	suite.useCase = use_cases.NewRefreshTokenUseCase(suite.authRepository)
}

func (suite *RefreshTokenUseCaseTestSuite) TestRefreshTokenWhenTokenIsValid() {
	token, refreshToken, _ := helpers.CreateTokens(validAuth.ID, validAuth.DomainArray(), validAuth.FlagsArray())
	suite.authRepository.On("Get", validAuth.ID).Return(validAuth, nil).Once()
	suite.authRepository.On("IsInvalidatedToken", refreshToken).Return(false).Once()
	arg := use_cases.RefreshTokenDTO{
		RefreshToken: refreshToken,
	}

	// Run the test
	// Need to sleep so the generated token is different
	time.Sleep(1 * time.Second)
	result, err := suite.useCase.RefreshToken(arg)

	// Verify the results
	suite.NoError(err)
	suite.NotEqualf(token, result.Token, "new token should be different")
	suite.NotEqualf(refreshToken, result.RefreshToken, "new refresh token should be different")
	newUserClaims, err := helpers.GetClaimsFromToken(result.Token)
	suite.NoError(err)
	suite.Equal(strconv.Itoa(int(validAuth.ID)), newUserClaims.ID)
	suite.authRepository.AssertExpectations(suite.T())
}

func (suite *RefreshTokenUseCaseTestSuite) TestRefreshTokenWhenTokenIsInvalid() {
	// Setup
	_, refreshToken, _ := helpers.CreateTokens(validAuth.ID, validAuth.DomainArray(), validAuth.FlagsArray())
	suite.authRepository.On("IsInvalidatedToken", refreshToken).Return(true).Once()
	arg := use_cases.RefreshTokenDTO{
		RefreshToken: refreshToken,
	}

	// Run the test
	_, err := suite.useCase.RefreshToken(arg)

	// Verify the results
	suite.Error(err, use_cases.ErrTokenInvalidated)
	suite.authRepository.AssertExpectations(suite.T())
	suite.authRepository.AssertNotCalled(suite.T(), "Get", validAuth.ID)
}

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

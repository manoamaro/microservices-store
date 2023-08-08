package use_cases

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/manoamaro/microservices-store/auth_service/internal/helpers"
	"github.com/manoamaro/microservices-store/auth_service/internal/use_cases"
	"github.com/manoamaro/microservices-store/auth_service/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type SignInUseCaseTestSuite struct {
	suite.Suite
	authRepository *test.MockAuthRepository
	useCase        use_cases.SignInUseCase
}

func TestSignInUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(SignInUseCaseTestSuite))
}

func (suite *SignInUseCaseTestSuite) SetupTest() {
	suite.authRepository = test.NewMockAuthRepository()
	suite.useCase = use_cases.NewSignInUseCase(suite.authRepository)
}

func (suite *SignInUseCaseTestSuite) TestSignInWithValidCredentials() {
	suite.authRepository.On("Authenticate", validAuth.Email, validAuth.Password).Return(&validAuth, true).Once()
	args := use_cases.SignInDTO{
		Email:         validAuth.Email,
		PlainPassword: validAuth.Password,
	}
	var domains jwt.ClaimStrings = validAuth.DomainArray()

	resultDTO, err := suite.useCase.SignIn(args)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resultDTO.Token)
	assert.NotEmpty(suite.T(), resultDTO.RefreshToken)
	userClaims, err := helpers.GetClaimsFromToken(resultDTO.Token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), strconv.Itoa(int(validAuth.ID)), userClaims.ID)
	assert.Equal(suite.T(), domains, userClaims.Audience)
	assert.Equal(suite.T(), validAuth.FlagsArray(), userClaims.Flags)
	suite.authRepository.AssertExpectations(suite.T())
}
